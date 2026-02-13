package kafka_sa

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/samber/lo"

	domainEvent "purchase/domain/event"
	"purchase/infra/idempotent"
	"purchase/infra/logx"
	"purchase/infra/mq"
	"purchase/infra/utils"
	"purchase/pkg/retry"
)

const (
	defaultRetryConsumerGroup = "retry_consumer_group"
	defaultDeadTopic          = "dead_topic"
	defaultMaxRetry           = 3
)

var retryTopic = map[time.Duration]string{
	5 * time.Second:  "retry_in_5s",
	10 * time.Second: "retry_in_10s",
	60 * time.Second: "retry_in_60s",
}

func retryBackoff(n int) time.Duration {
	if n == 1 {
		return 5 * time.Second
	}
	if n == 2 {
		return 10 * time.Second
	}
	return 60 * time.Second
}

type kafkaSubscriber struct {
	client sarama.Client
	// event is topic, handler is consumer group
	handlers map[domainEvent.Event][]domainEvent.Handler
	// address   []string
	idp       idempotent.Idempotent
	consumers []*Consumer
	pub       mq.Publisher
	rlq       *retryRouter
	dlq       *dlqRouter
	conf      *sarama.Config
}

func NewKafkaSubscriber(cfg config.Config, idp idempotent.Idempotent, handlerAgg domainEvent.HandlerAggregator, pub mq.Publisher) (mq.Subscriber, error) {
	address := make([]string, 0)
	err := cfg.Value("kafka.address").Scan(&address)
	if err != nil {
		return nil, err
	}

	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true               // 成功发送的消息将写到 Successes 通道
	conf.Consumer.Return.Errors = true                  // 消费时错误信息将写到 Errors 通道
	conf.Consumer.Fetch.Default = 3 * 1024 * 1024       // 默认请求的字节数
	conf.Consumer.Offsets.Initial = sarama.OffsetNewest // 从最新的 offset 读取，如果设置为 OffsetOldest 则从最旧的 offset 读取
	conf.Consumer.Offsets.AutoCommit.Enable = false     // 不自动提交offset

	client, err := sarama.NewClient(address, conf)
	if err != nil {
		return nil, err
	}

	s := &kafkaSubscriber{
		client:   client,
		handlers: handlerAgg.Build(),
		// address:  address,
		idp:  idp,
		pub:  pub,
		conf: conf,
	}
	return s, nil
}

func (s *kafkaSubscriber) Subscribe(ctx context.Context) {
	transferHandler := func(events []domainEvent.Event, h domainEvent.Handler) mq.Handler {
		return func(ctx context.Context, m *mq.Message) error {
			for _, e := range events {
				if e.EventName() == m.EventName() {
					evt, err := e.Decode(m.Body)
					if err != nil {
						return err
					}
					return h(ctx, evt)
				}
			}
			return nil
		}
	}

	handlerEvents := make(map[string][]domainEvent.Event)
	handlerMap := make(map[string]domainEvent.Handler)
	for e, handlers := range s.handlers {
		for _, h := range handlers {
			handlerName := utils.GetMethodName(h)
			handlerEvents[handlerName] = append(handlerEvents[handlerName], e)
			handlerMap[handlerName] = h
		}
	}

	for handlerName, events := range handlerEvents {
		topics := lo.Map(events, func(e domainEvent.Event, _ int) string {
			return e.EventName()
		})
		c, err := NewConsumer(s, topics, handlerName, transferHandler(events, handlerMap[handlerName]), false)
		if err != nil {
			logx.Error(ctx, "new consume err", slog.Any("error", err))
			continue
		}
		s.registerConsumer(c)
		go c.Run(ctx)
	}

	s.rlq = newRetryRouter(s.pub)

	s.dlq = newDlqRouter(s.pub)

	s.consumeRetryTopic(ctx)
}

func (s *kafkaSubscriber) registerConsumer(c *Consumer) {
	s.consumers = append(s.consumers, c)
}

func (s *kafkaSubscriber) Close() error {
	for _, c := range s.consumers {
		c.Close()
	}
	return nil
}

func NewConsumer(sub *kafkaSubscriber, topics []string, consumerGroup string, handler mq.Handler, isConsumeRlq bool) (*Consumer, error) {
	cg, err := sarama.NewConsumerGroupFromClient(consumerGroup, sub.client)
	if err != nil {
		return nil, err
	}

	c := &Consumer{
		sub:          sub,
		handler:      handler,
		isConsumeRlq: isConsumeRlq,
		cg:           cg,
		topics:       topics,
		ctx:          context.Background(),
	}

	return c, nil
}

// Consumer is a consumer group
type Consumer struct {
	sub          *kafkaSubscriber
	handler      mq.Handler
	isConsumeRlq bool
	ctx          context.Context
	cg           sarama.ConsumerGroup
	topics       []string
}

func (c *Consumer) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logx.Info(ctx, "consumer stopped due to context cancellation", slog.Any("topics", c.topics))
			return
		default:
		}
		err := c.cg.Consume(ctx, c.topics, c) // 传入定义好的 ConsumerGroupHandler 结构体
		if err != nil {
			logx.Errorf(ctx, "consume error: %#v\n", err)
		}
	}
}

// handleMessage 幂等消费消息
func (c *Consumer) handleMessage(ctx context.Context, sess sarama.ConsumerGroupSession, m *mq.Message, h mq.Handler) {
	ok, err := c.sub.idp.SetKeyPendingWithDDL(ctx, m.ID, time.Second*time.Duration(10*60))
	if err != nil {
		logx.Error(ctx, "subscriber SetKeyPendingWithDDL fail", slog.Any("error", err), slog.String("key", m.ID), slog.Any("message", m))
		return
	}
	// 已经消费过
	if !ok {
		state, err := c.sub.idp.GetKeyState(ctx, m.ID)
		if err != nil {
			logx.Error(ctx, "subscriber GetKeyState fail", slog.Any("error", err), slog.String("key", m.ID), slog.Any("message", m))
			return
		}
		if state == idempotent.Done {
			// 已经消费过，且消费成功了, 即使提交失败了，只要后面的message能成功，之前的message都会被提交
			if err = c.commitMessage(ctx, sess, m); err != nil {
				logx.Error(ctx, "failed to commit messages:", slog.Any("error", err), slog.Any("message", m))
			}
			return
		}
		// 消费重试队列里的消息，只是将消息推送到初始topic，理论不会失败，不要再往重试队列里投递，避免复杂化
		if !c.isConsumeRlq {
			// 已经消费过，但消费失败了
			c.ReconsumeLater(sess, m)
		}
		return
	}
	// 未消费过，执行消费逻辑
	err = h(ctx, m)
	// 消费失败
	if err != nil {
		logx.Error(ctx, "consume message fail", slog.Any("error", err), slog.Any("message", m))
		// 消费失败, 清除掉幂等key
		// 如果RemoveFailKey失败，ReconsumeLater成功，重会等到key过期删除，然后被消费
		// 如果RemoveFailKey成功，ReconsumeLater失败，队列未提交ack（下一次offset提交成功，会把之前的offset一起提交），不会再次消费，需要人工介入resume
		// 如果RemoveFailKey失败，ReconsumeLater失败，队列未提交ack，会等到key过期删除，不会再次被消费, 需要人工介入resume
		err = retry.Run(func() error {
			return c.sub.idp.RemoveFailKey(ctx, m.ID)
		}, 2)
		if err != nil {
			logx.Error(ctx, "RemoveFailKey fail after retry 2 times", slog.Any("error", err), slog.String("key", m.ID))
		}
		if !c.isConsumeRlq {
			c.ReconsumeLater(sess, m)
		}
		return
	}
	// 消费成功
	err = retry.Run(func() error {
		return c.sub.idp.UpdateKeyDone(ctx, m.ID)
	}, 2)
	if err != nil {
		logx.Error(ctx, "UpdateKeyDone fail after retry 2 times", slog.Any("error", err), slog.String("key", m.ID))
	}
	err = retry.Run(func() error {
		return c.commitMessage(ctx, sess, m)
	}, 2)
	if err != nil {
		logx.Error(ctx, "consumer ack fail after retry 2 times", slog.Any("error", err), slog.Any("message", m))
	}
	return
}

func (c *Consumer) buildMessage(m *sarama.ConsumerMessage) (*mq.Message, error) {
	msg := &mq.Message{
		Body: m.Value,
	}
	msg.SetPartition(m.Partition)
	for _, header := range m.Headers {
		if string(header.Key) == mq.MessageID {
			msg.ID = string(header.Value)
		}
		if string(header.Key) == HeaderPropertyKey {
			props := &mq.Header{}
			err := json.Unmarshal(header.Value, props)
			if err != nil {
				return nil, err
			}
			msg.SetHeader(*props)
		}
	}
	msg.SetRawMessage(m)

	return msg, nil
}

// commitMessage 使用指定的 session 提交消息，避免并发竞争
func (c *Consumer) commitMessage(ctx context.Context, sess sarama.ConsumerGroupSession, m *mq.Message) error {
	raw, ok := m.RawMessage().(*sarama.ConsumerMessage)
	if !ok {
		return fmt.Errorf("invalid raw message type: %T", m.RawMessage())
	}
	sess.MarkMessage(raw, "")
	sess.Commit()
	return nil
}

func (c *Consumer) Close() error {
	err := c.cg.Close()
	fmt.Print("close kafka reader")
	if err != nil {
		return err
	}
	return nil
}

func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		m, err := c.buildMessage(msg)
		fmt.Println(fmt.Sprintf("收到消息topic=%s, partition=%d, offset=%d", msg.Topic, msg.Partition, msg.Offset))
		if err != nil {
			logx.Error(c.ctx, "buildMessage fail", slog.Any("error", err))
			continue
		}
		if c.isConsumeRlq {
			if err := c.redelivery(c.ctx, sess, m); err != nil {
				logx.Error(c.ctx, "redelivery fail", slog.Any("error", err))
			}
		} else {
			c.handleMessage(c.ctx, sess, m, c.handler)
		}
	}
	return nil
}

// consumeRetryTopic 重新投递到原来的队列
func (s *kafkaSubscriber) consumeRetryTopic(ctx context.Context) {
	retryTopics := lo.MapToSlice(retryTopic, func(k time.Duration, v string) string {
		return v
	})
	// 从重试队列拉去消息，发送到消息原来的topic重新消费
	go func() {
		c, err := NewConsumer(s, retryTopics, defaultRetryConsumerGroup, nil, true)
		if err != nil {
			logx.Error(ctx, "new consume err", slog.Any("error", err))
			return
		}
		s.registerConsumer(c)
		c.Run(ctx)
	}()
}

func (c *Consumer) redelivery(ctx context.Context, sess sarama.ConsumerGroupSession, m *mq.Message) error {
	// 未到消费时间，sleep, 此处应该调用kafka的pause和resume api，不然会导致重平衡，但是此kafka client不支持
	now := time.Now()
	if expTime := m.DeliveryTime().Add(m.DelayTime()); expTime.After(now) {
		partitions := map[string][]int32{
			m.RetryTopic(): {m.Partition()},
		}
		c.cg.Pause(partitions)
		fmt.Println(fmt.Sprintf("bizCODE: %s, TOPIC:%s,  PARTITION:%d, 开始pause: %v, 持续%v ", m.BizCode(), m.RetryTopic(), m.Partition(), time.Now(), expTime.Sub(now)))
		time.Sleep(expTime.Sub(now))
		c.cg.Resume(partitions)
		fmt.Println(fmt.Sprintf("bizCODE: %s, TOPIC:%s,  PARTITION:%d, resume: %v", m.BizCode(), m.RetryTopic(), m.Partition(), time.Now()))
	}
	// 清空死信topic和重投topic
	// 让消息发送到原始topic里
	m.SetDeadTopic("")
	m.SetRetryTopic("")
	fmt.Println(fmt.Sprintf("bizCODE: %s发送消息到topic:%s", m.BizCode(), m.EventName()))
	err := c.sub.pub.Publish(ctx, m)
	if err != nil {
		logx.Error(ctx, fmt.Sprintf("retry message ,bizCODE: %s发送消息到topic:%s", m.BizCode(), m.EventName()), slog.Any("error", err))
		return err
	}
	err = c.commitMessage(ctx, sess, m)
	if err != nil {
		logx.Error(ctx, "retry message consumer ack fail", slog.Any("error", err))
	}
	return err
}

func (c *Consumer) ReconsumeLater(sess sarama.ConsumerGroupSession, m *mq.Message) {
	// m.SetOriginalTopic(c.topic)
	m.IncrReconsumeTimes()
	m.SetDelayTime(retryBackoff(m.ReconsumeTimes()))
	retryMessage := RetryMessage{
		Message:  m,
		consumer: c,
		sess:     sess,
	}
	if m.ReconsumeTimes() > c.sub.dlq.maxRetry() {
		c.sub.dlq.Chan() <- retryMessage
	} else {
		c.sub.rlq.Chan() <- retryMessage
	}
}

type RetryMessage struct {
	*mq.Message
	consumer *Consumer
	sess     sarama.ConsumerGroupSession
}

func (r RetryMessage) Commit(ctx context.Context) error {
	return r.consumer.commitMessage(ctx, r.sess, r.Message)
}

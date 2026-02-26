package kafka_sa

import (
	"context"
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

func genRetryTopic(d time.Duration) string {
	if t, ok := retryTopic[d]; ok {
		return t
	}
	return retryTopic[60*time.Second]
}

type kafkaSubscriber struct {
	client sarama.Client
	// event is topic, handler is consumer group
	handlers  map[domainEvent.Event][]domainEvent.Handler
	idp       idempotent.Idempotent
	consumers []*Consumer
	pub       mq.Publisher
	rlq       *messageRouter
	dlq       *messageRouter
	conf      *sarama.Config
	marshaler MessageMarshaler
}

func NewKafkaSubscriber(cfg config.Config, idp idempotent.Idempotent, handlerAgg domainEvent.HandlerAggregator, pub mq.Publisher, idg mq.IDGenFunc) (mq.Subscriber, error) {
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
		client:    client,
		handlers:  handlerAgg.Build(),
		idp:       idp,
		pub:       pub,
		conf:      conf,
		marshaler: NewMessageMarshaler(idg),
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

	s.rlq = newMessageRouter("retryRouter", s.pub, func(msg *RetryMessage) {
		rtyTopic := genRetryTopic(msg.DelayTime())
		msg.SetRetryTopic(rtyTopic)
	})

	s.dlq = newMessageRouter("dlqRouter", s.pub, func(msg *RetryMessage) {
		msg.SetDeadTopic(defaultDeadTopic)
	})

	s.consumeRetryTopic(ctx)
}

func (s *kafkaSubscriber) registerConsumer(c *Consumer) {
	s.consumers = append(s.consumers, c)
}

func (s *kafkaSubscriber) Close() error {
	if s.rlq != nil {
		s.rlq.close()
	}
	if s.dlq != nil {
		s.dlq.close()
	}
	for _, c := range s.consumers {
		c.Close()
	}
	return s.client.Close()
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

// idempotentBackoff 幂等操作失败时的退避时间，避免热循环
const idempotentBackoff = 3 * time.Second

// handleMessage 幂等消费消息
func (c *Consumer) handleMessage(ctx context.Context, sess sarama.ConsumerGroupSession, m *mq.Message, h mq.Handler) {
	ok, err := c.sub.idp.SetKeyPendingWithDDL(ctx, m.ID, time.Second*time.Duration(10*60))
	if err != nil {
		logx.Error(ctx, "subscriber SetKeyPendingWithDDL fail", slog.Any("error", err), slog.String("key", m.ID), slog.Any("message", m))
		// backoff 避免热循环：Redis/MySQL 短暂不可用时，不 commit 让 Kafka 重新投递，但需要降低重试频率
		time.Sleep(idempotentBackoff)
		return
	}
	// 已经消费过（key 已存在）
	if !ok {
		state, err := c.sub.idp.GetKeyState(ctx, m.ID)
		if err != nil {
			logx.Error(ctx, "subscriber GetKeyState fail", slog.Any("error", err), slog.String("key", m.ID), slog.Any("message", m))
			// backoff 避免热循环
			time.Sleep(idempotentBackoff)
			return
		}
		switch state {
		case idempotent.Done:
			// 已经消费过，且消费成功了, 即使提交失败了，只要后面的 message 能成功，之前的 message 都会被提交
			if err = c.commitMessage(ctx, sess, m); err != nil {
				logx.Error(ctx, "failed to commit messages", slog.Any("error", err), slog.String("key", m.ID), slog.Any("message", m))
			}
		case idempotent.Pending:
			// 有另一个消费者正在处理这条消息，不做任何操作
			// 不 commit：让 Kafka rebalance 后重新投递
			// 不 ReconsumeLater：避免与正在处理的消费者产生重复消费
			logx.Info(ctx, "message is being consumed by another consumer, skip",
				slog.String("key", m.ID), slog.String("bizCode", m.BizCode()))
		default:
			// 未知状态，走重试逻辑
			logx.Error(ctx, "unknown idempotent state", slog.String("state", state), slog.String("key", m.ID))
			if !c.isConsumeRlq {
				c.ReconsumeLater(sess, m)
			}
		}
		return
	}
	// 未消费过，执行消费逻辑
	err = h(ctx, m)
	// 消费失败
	if err != nil {
		logx.Error(ctx, "consume message fail", slog.Any("error", err), slog.Any("message", m))
		// 消费失败, 清除掉幂等 key
		// 如果 RemoveFailKey 失败，ReconsumeLater 成功，会等到 key 过期删除，然后被消费
		// 如果 RemoveFailKey 成功，ReconsumeLater 失败，队列未提交 ack（下一次 offset 提交成功，会把之前的 offset 一起提交），不会再次消费，需要人工介入
		// 如果 RemoveFailKey 失败，ReconsumeLater 失败，队列未提交 ack，会等到 key 过期删除，不会再次被消费，需要人工介入
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
		// UpdateKeyDone 失败，key 状态仍为 Pending
		// 由于上面 Pending 分支不会触发 ReconsumeLater，即使消息被重新投递也只会 skip，不会重复消费
		logx.Error(ctx, "UpdateKeyDone fail after retry 2 times, message consumed but state not updated",
			slog.Any("error", err), slog.String("key", m.ID))
	}
	err = retry.Run(func() error {
		return c.commitMessage(ctx, sess, m)
	}, 2)
	if err != nil {
		logx.Error(ctx, "consumer ack fail after retry 2 times", slog.Any("error", err), slog.Any("message", m))
	}
}

func (c *Consumer) buildMessage(ctx context.Context, m *sarama.ConsumerMessage) (*mq.Message, error) {
	return c.sub.marshaler.Unmarshal(ctx, m)
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
	logx.Info(c.ctx, "closing kafka consumer", slog.Any("topics", c.topics))
	return c.cg.Close()
}

func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		logx.Info(c.ctx, "received message", slog.String("topic", msg.Topic), slog.Int("partition", int(msg.Partition)), slog.Int64("offset", msg.Offset))
		m, err := c.buildMessage(c.ctx, msg)
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
	// 未到消费时间，pause partition 并等待
	now := time.Now()
	// Sarama 为每个 topic-partition 分配独立的 goroutine
	// 执行 ConsumeClaim，Pause 也是按 topic-partition 粒度生效的。
	// retry_in_60s 的消息在等待时，retry_in_5s 的消费完全不受影响。
	if expTime := m.DeliveryTime().Add(m.DelayTime()); expTime.After(now) {
		partitions := map[string][]int32{
			m.RetryTopic(): {m.Partition()},
		}
		// Pause 只是停止 fetch，不影响心跳。
		// Sarama 的心跳是独立 goroutine 发送的
		c.cg.Pause(partitions)
		waitDuration := expTime.Sub(now)
		logx.Info(ctx, "pause partition for redelivery",
			slog.String("bizCode", m.BizCode()),
			slog.String("retryTopic", m.RetryTopic()),
			slog.Int("partition", int(m.Partition())),
			slog.Duration("waitDuration", waitDuration),
		)
		select {
		case <-time.After(waitDuration):
		case <-ctx.Done():
			c.cg.Resume(partitions)
			return ctx.Err()
		}
		c.cg.Resume(partitions)
		logx.Info(ctx, "resume partition after redelivery wait",
			slog.String("bizCode", m.BizCode()),
			slog.String("retryTopic", m.RetryTopic()),
			slog.Int("partition", int(m.Partition())),
		)
	}
	// 清空死信topic和重投topic，让消息发送到原始topic里
	m.SetDeadTopic("")
	m.SetRetryTopic("")
	logx.Info(ctx, "redelivery message to original topic",
		slog.String("bizCode", m.BizCode()),
		slog.String("topic", m.EventName()),
	)
	err := retry.Run(func() error {
		return c.sub.pub.Publish(ctx, m)
	}, 3)
	if err != nil {
		logx.Error(ctx, "redelivery publish fail after retries",
			slog.Any("error", err),
			slog.String("bizCode", m.BizCode()),
			slog.String("topic", m.EventName()),
			slog.String("messageID", m.ID),
		)
		return err
	}
	err = retry.Run(func() error {
		return c.commitMessage(ctx, sess, m)
	}, 2)
	if err != nil {
		logx.Error(ctx, "redelivery commit fail after retries", slog.Any("error", err), slog.String("messageID", m.ID))
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

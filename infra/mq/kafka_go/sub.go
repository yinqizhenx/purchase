package kafka_go

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/segmentio/kafka-go"

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
	handlers  map[domainEvent.Event][]domainEvent.Handler
	address   []string
	idp       idempotent.Idempotent
	consumers []*Consumer
	pub       mq.Publisher
	rlq       *retryRouter
	dlq       *dlqRouter
}

func NewKafkaSubscriber(cfg config.Config, idp idempotent.Idempotent, handlerAgg domainEvent.HandlerAggregator, pub mq.Publisher) (mq.Subscriber, error) {
	address := make([]string, 0)
	err := cfg.Value("kafka.address").Scan(&address)
	if err != nil {
		return nil, err
	}
	s := &kafkaSubscriber{
		handlers: handlerAgg.Build(),
		address:  address,
		idp:      idp,
		pub:      pub,
	}
	return s, nil
}

func (s *kafkaSubscriber) Subscribe(ctx context.Context) {
	transferHandler := func(e domainEvent.Event, h domainEvent.Handler) mq.Handler {
		return func(ctx context.Context, m *mq.Message) error {
			if e.EventName() == m.EventName() {
				msg, err := e.Decode(m.Body)
				if err != nil {
					return err
				}
				return h(ctx, msg)
			}
			return nil
		}
	}

	for e, handlers := range s.handlers {
		for _, h := range handlers {
			c := NewConsumer(s, e.EventName(), utils.GetMethodName(h), transferHandler(e, h), false)
			s.registerConsumer(c)
			go c.Run(ctx)
		}
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

func NewConsumer(sub *kafkaSubscriber, topic string, consumerGroup string, handler mq.Handler, isConsumeRlq bool) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  sub.address,
		GroupID:  consumerGroup,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	c := &Consumer{
		sub:          sub,
		reader:       reader,
		handler:      handler,
		isConsumeRlq: isConsumeRlq,
	}

	if c.isConsumeRlq {
		c.handler = c.redelivery
	}

	return c
}

type Consumer struct {
	sub          *kafkaSubscriber
	reader       *kafka.Reader
	handler      mq.Handler
	isConsumeRlq bool
}

func (c *Consumer) Run(ctx context.Context) {
	for {
		kmsg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if err != io.EOF {
				logx.Error(ctx, "kafka fetch message fail", slog.Any("error", err))
			}
			continue
		}
		m, err := c.buildMessage(kmsg)
		if err != nil {
			logx.Error(ctx, "buildMessage fail", slog.Any("error", err))
			continue
		}
		c.handleMessage(ctx, m, c.handler)
	}
}

// handleMessage 幂等消费消息
func (c *Consumer) handleMessage(ctx context.Context, m *mq.Message, h mq.Handler) {
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
			if err = c.CommitMessage(ctx, m); err != nil {
				logx.Error(ctx, "failed to commit messages:", slog.Any("error", err), slog.Any("message", m))
			}
		} else {
			if !c.isConsumeRlq {
				// 已经消费过，但消费失败了
				c.ReconsumeLater(m)
			}
		}
	} else {
		// 未消费过，执行消费逻辑
		err = h(ctx, m)
		if err != nil {
			logx.Error(ctx, "consume message fail", slog.Any("error", err), slog.Any("message", m))
			// 消费失败, 清除掉幂等key
			// 如果RemoveFailKey失败，ReconsumeLater成功，重会等到key过期删除，然后被消费
			// 如果RemoveFailKey成功，ReconsumeLater失败，不会再次消费
			// 如果RemoveFailKey失败，ReconsumeLater失败，队列未提交ack，会等到key过期删除，不会再次被消费
			err = retry.Run(func() error {
				return c.sub.idp.RemoveFailKey(ctx, m.ID)
			}, 2)
			if err != nil {
				logx.Error(ctx, "RemoveFailKey fail after retry 2 times", slog.Any("error", err), slog.String("key", m.ID))
			}
			if !c.isConsumeRlq {
				c.ReconsumeLater(m)
			}
		} else {
			// 消费采成功
			err = retry.Run(func() error {
				return c.sub.idp.UpdateKeyDone(ctx, m.ID)
			}, 2)
			if err != nil {
				logx.Error(ctx, "UpdateKeyDone fail after retry 2 times", slog.Any("error", err), slog.String("key", m.ID))
			}
			err = retry.Run(func() error {
				return c.CommitMessage(ctx, m)
			}, 2)
			if err != nil {
				logx.Error(ctx, "consumer ack fail after retry 2 times", slog.Any("error", err), slog.Any("message", m))
			}
		}
	}
}

func (c *Consumer) buildMessage(m kafka.Message) (*mq.Message, error) {
	msg := &mq.Message{
		Body: m.Value,
	}
	for _, header := range m.Headers {
		if header.Key == mq.MessageID {
			msg.ID = string(header.Value)
		}
		if header.Key == HeaderPropertyKey {
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

func (c *Consumer) CommitMessage(ctx context.Context, m *mq.Message) error {
	return c.reader.CommitMessages(ctx, m.RawMessage().(kafka.Message))
}

func (c *Consumer) Close() error {
	err := c.reader.Close()
	fmt.Print("close kafka reader")
	if err != nil {
		return err
	}
	return nil
}

// consumeRetryTopic 重新投递到原来的队列
func (s *kafkaSubscriber) consumeRetryTopic(ctx context.Context) {
	for _, topic := range retryTopic {
		// 从重试队列拉去消息，发送到消息原来的topic重新消费
		go func(t string) {
			c := NewConsumer(s, t, defaultRetryConsumerGroup, nil, true)
			c.Run(ctx)
		}(topic)
	}
}

func (c *Consumer) redelivery(ctx context.Context, m *mq.Message) error {
	// 清空死信topic和重投topic
	// 让消息发送到原始topic里
	m.SetDeadTopic("")
	m.SetRetryTopic("")
	// todo 时间还未到，需要不投
	// // 未到消费时间，sleep, 此处应该调用kafka的pause和resume api，不然会导致重平衡，但是此kafka client不支持
	// if expTime := m.RedeliveryTime().Add(m.DelayTime()); expTime.After(time.Now()) {
	// 	fmt.Println("开始sleep ", time.Now().Sub(expTime).Seconds())
	// 	time.Sleep(time.Now().Sub(expTime))
	// }
	err := c.sub.pub.Publish(ctx, m)
	if err != nil {
		return err
	}
	err = c.CommitMessage(ctx, m)
	if err != nil {
		logx.Error(ctx, "retry message consumer ack fail", slog.Any("error", err))
	}
	return err
}

func (c *Consumer) ReconsumeLater(m *mq.Message) {
	// m.SetOriginalTopic(c.topic)
	m.IncrReconsumeTimes()
	m.SetDelayTime(retryBackoff(m.ReconsumeTimes()))
	retryMessage := RetryMessage{
		Message:          m,
		MessageCommitter: c,
	}
	if m.ReconsumeTimes() > c.sub.dlq.maxRetry() {
		c.sub.dlq.Chan() <- retryMessage
	} else {
		c.sub.rlq.Chan() <- retryMessage
	}
}

type RetryMessage struct {
	*mq.Message
	mq.MessageCommitter
}

func (r RetryMessage) Commit(ctx context.Context) error {
	return r.MessageCommitter.CommitMessage(ctx, r.Message)
}

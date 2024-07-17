package kafka

import (
	"context"
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
	defaultRetryTopic = "KAFKA_RETRY_TOPIC"
	defaultDeadTopic  = "KAFKA_DEAD_TOPIC"
	defaultMaxRetry   = 3
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
	handlers map[domainEvent.Event][]domainEvent.Handler
	address  []string
	idp      idempotent.Idempotent
}

func NewKafkaSub(cfg config.Config, idp idempotent.Idempotent, handlerAgg domainEvent.HandlerAggregator) (mq.Subscriber, error) {
	address := make([]string, 0)
	err := cfg.Value("kafka.address").Scan(&address)
	if err != nil {
		return nil, err
	}
	s := &kafkaSubscriber{
		handlers: handlerAgg.Build(),
		address:  address,
		idp:      idp,
	}
	return s, nil
}

func (s *kafkaSubscriber) Subscribe(ctx context.Context) {
	transferHandler := func(e domainEvent.Event, h domainEvent.Handler) mq.Handler {
		return func(ctx context.Context, m *mq.Message) error {
			if e.EventName() == m.HeaderGet(mq.EventName) {
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
			c, err := NewConsumer(e.EventName(), utils.GetMethodName(h), s.address, transferHandler(e, h), s.idp)
			if err != nil {
				logx.Errorf(ctx, "new consuer error")
				continue
			}
			go c.Run(ctx)
		}
	}
}

func (s *kafkaSubscriber) Close() error {
	return nil
}

func NewConsumer(topic string, consumerGroup string, address []string, handler mq.Handler, idp idempotent.Idempotent) (*Consumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  address,
		GroupID:  consumerGroup,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	rlqPolicy := &RLQPolicy{
		RetryLetterTopic: defaultRetryTopic,
		Address:          address,
		GroupID:          "group-b",
	}
	retryRouter, err := newRetryRouter(rlqPolicy, reader)
	if err != nil {
		return nil, err
	}

	dlqPolicy := &DLQPolicy{
		MaxDeliveries:   3,
		DeadLetterTopic: defaultDeadTopic,
		Address:         address,
		GroupID:         "group-b",
	}
	dlqRouter, err := newDlqRouter(dlqPolicy, reader)
	if err != nil {
		return nil, err
	}

	c := &Consumer{
		reader:  reader,
		handler: handler,
		topic:   topic,
		idp:     idp,
		rlq:     retryRouter,
		dlq:     dlqRouter,
	}
	go c.consumeRetryTopic(context.Background(), address)
	return c, nil
}

type Consumer struct {
	reader       *kafka.Reader
	handler      mq.Handler
	idp          idempotent.Idempotent
	retryEnabled bool
	rlq          *retryRouter
	dlq          *dlqRouter
	topic        string
}

func (c *Consumer) Run(ctx context.Context) {
	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if err != io.EOF {
				logx.Error(ctx, "kafka fetch message fail", slog.Any("error", err))
			}
			continue
		}
		c.handleMessage(ctx, m, c.handler)
	}
}

// handleMessage 幂等消费消息
func (c *Consumer) handleMessage(ctx context.Context, m kafka.Message, h mq.Handler) {
	msg := NewMessage(&m)
	key := msg.PropsMessageID()
	ok, err := c.idp.SetKeyPendingWithDDL(ctx, key, time.Second*time.Duration(10*60))
	if err != nil {
		logx.Error(ctx, "subscriber SetKeyPendingWithDDL fail", slog.Any("error", err), slog.String("key", key), slog.Any("message", m))
		return
	}
	// 已经消费过
	if !ok {
		state, err := c.idp.GetKeyState(ctx, key)
		if err != nil {
			logx.Error(ctx, "subscriber GetKeyState fail", slog.Any("error", err), slog.String("key", key), slog.Any("message", m))
			return
		}
		if state == idempotent.Done {
			// 已经消费过，且消费成功了, 即使提交失败了，只要后面的message能成功，之前的message都会被提交
			if err = c.reader.CommitMessages(ctx, m); err != nil {
				logx.Error(ctx, "failed to commit messages:", slog.Any("error", err), slog.Any("message", m))
			}
		} else {
			// 已经消费过，但消费失败了
			c.ReconsumeLater(msg)
		}
	} else {
		// 未消费过，执行消费逻辑
		message := &mq.Message{
			Body: m.Value,
		}
		message.HeaderSet(mq.EventName, msg.PropsEventName())
		err = h(ctx, message)
		if err != nil {
			logx.Error(ctx, "consume message fail", slog.Any("error", err), slog.Any("message", m))
			// 消费失败, 清除掉幂等key
			// 如果RemoveFailKey失败，ReconsumeLater成功，重会等到key过期删除，然后被消费
			// 如果RemoveFailKey成功，ReconsumeLater失败，不会再次消费
			// 如果RemoveFailKey失败，ReconsumeLater失败，队列未提交ack，会等到key过期删除，不会再次被消费
			err = retry.Run(func() error {
				return c.idp.RemoveFailKey(ctx, key)
			}, 2)
			if err != nil {
				logx.Error(ctx, "RemoveFailKey fail after retry 2 times", slog.Any("error", err), slog.String("key", key))
			}
			c.ReconsumeLater(msg)
		} else {
			// 消费采成功
			err = retry.Run(func() error {
				return c.idp.UpdateKeyDone(ctx, key)
			}, 2)
			if err != nil {
				logx.Error(ctx, "UpdateKeyDone fail after retry 2 times", slog.Any("error", err), slog.String("key", key))
			}
			err = retry.Run(func() error {
				return c.reader.CommitMessages(ctx, m)
			}, 2)
			if err != nil {
				logx.Error(ctx, "consumer ack fail after retry 2 times", slog.Any("error", err), slog.Any("message", m))
			}
		}
	}
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
func (c *Consumer) consumeRetryTopic(ctx context.Context, address []string) {
	w := &kafka.Writer{
		// Topic:    topic,
		Addr:     kafka.TCP(address...),
		Balancer: &kafka.LeastBytes{},
	}
	for _, topic := range retryTopic {
		// 从重试队列拉去消息，发送到消息原来的topic重新消费
		go func(t string) {
			r := kafka.NewReader(kafka.ReaderConfig{
				Brokers:  address,
				GroupID:  "group-retry",
				Topic:    t,
				MinBytes: 10e3, // 10KB
				MaxBytes: 10e6, // 10MB
			})
			c.redelivery(ctx, r, w)
		}(topic)
	}
}

func (c *Consumer) redelivery(ctx context.Context, r *kafka.Reader, w *kafka.Writer) {
	for {
		m, err := r.FetchMessage(context.Background())
		if err != nil {
			if err != io.EOF {
				logx.Error(ctx, "kafka fetch message fail", slog.Any("error", err))
			}
			continue
		}
		msg := NewMessage(&m)
		msg.SetMessageTopic(msg.PropsRealTopic())
		// // 未到消费时间，sleep, 此处应该调用kafka的pause和resume api，不然会导致重平衡，但是此kafka client不支持
		// if expTime := msg.PropsRedeliveryTime().Add(msg.PropsDelayTime()); expTime.After(time.Now()) {
		// 	fmt.Println("开始sleep ", time.Now().Sub(expTime).Seconds())
		// 	time.Sleep(time.Now().Sub(expTime))
		// }

		// 未消费过，执行消费逻辑
		kmsg, err := msg.ToKafkaMessage()
		if err != nil {
			logx.Error(ctx, "message transfer to kafka message fail", slog.Any("message", msg), slog.Any("error", err))
			continue
		}
		err = retry.Run(func() error { return w.WriteMessages(ctx, *kmsg) }, 2)
		if err != nil {
			logx.Error(ctx, "write consume message to topic fail after retry 2 times", slog.String("topic", msg.PropsRealTopic()), slog.Any("error", err))
			continue
		}

		kmsg.Topic = msg.PropsRetryTopic()
		err = r.CommitMessages(ctx, *kmsg)
		if err != nil {
			logx.Error(ctx, "retry message consumer ack fail", slog.Any("error", err))
		}
	}
}

func (c *Consumer) ReconsumeLater(m *Message) {
	m.SetRealTopic(c.topic)
	m.IncrReconsumeTimes()
	m.SetDelayTime(retryBackoff(m.PropsReconsumeTimes()))
	if m.PropsReconsumeTimes() > c.dlq.maxRetry() {
		c.dlq.Chan() <- m
	} else {
		c.rlq.Chan() <- m
	}
}

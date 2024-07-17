package kafka

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	// "log"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/segmentio/kafka-go"

	"purchase/infra/idempotent"
	"purchase/infra/logx"
	"purchase/infra/mq"
	"purchase/pkg/retry"
)

var (
	_ mq.Publisher = (*kafkaPublisher)(nil)
	// _ mq.Subscriber = (*kafkaSubscriber)(nil)
	// _ Event      = (*Message)(nil)
)

var ProviderSet = wire.NewSet(NewKafkaPublisher, NewKafkaSub)

type kafkaPublisher struct {
	writer *kafka.Writer
	topic  string
	idg    mq.IDGenFunc
}

func (s *kafkaPublisher) Publish(ctx context.Context, msg *mq.Message) error {
	kMsg := &kafka.Message{
		Key:   []byte(msg.HeaderGet(mq.EntityID)),
		Value: msg.Body,
	}
	m := NewMessage(kMsg)
	m.SetEventName(msg.HeaderGet(mq.EventName))
	m.SetMessageID(s.idg())

	kMsg, err := m.ToKafkaMessage()
	if err != nil {
		return err
	}
	err = s.writer.WriteMessages(ctx, *kMsg)
	if err != nil {
		return err
	}
	return nil
}

func (s *kafkaPublisher) Close() error {
	err := s.writer.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewKafkaPublisher(c config.Config, idg mq.IDGenFunc) (mq.Publisher, error) {
	address := make([]string, 0)
	err := c.Value("kafka.address").Scan(&address)
	if err != nil {
		return nil, err
	}
	topic, err := c.Value("kafka.topic").String()
	if err != nil {
		return nil, err
	}
	w := &kafka.Writer{
		Topic:    topic,
		Addr:     kafka.TCP(address...),
		Balancer: &kafka.LeastBytes{},
	}
	return &kafkaPublisher{writer: w, topic: topic, idg: idg}, nil
}

type kafkaSubscriber struct {
	reader       *kafka.Reader
	topic        string
	idp          idempotent.Idempotent
	retryEnabled bool
	rlq          *retryRouter
	dlq          *dlqRouter
	h            mq.Handler
}

type Option struct {
	topic        string
	retryEnabled bool
}

type retryPolicy struct {
	rlqTopic    string
	dlqTopic    string
	backOff     retry.BackoffPolicy
	rlqReader   *kafka.Reader // 读取重试队列
	retryWriter *kafka.Writer // 发送到重试队列
	deadWriter  *kafka.Writer
}

// type DLQ struct {
// 	maxRetry int
// 	topic    string
// 	writer   *kafka.Writer // 发送到死信队列
// }

const (
	defaultRetryTopic = "KAFKA_RETRY_TOPIC"
	defaultDeadTopic  = "KAFKA_DEAD_TOPIC"
	defaultMaxRetry   = 3
)

func NewKafkaSubscriber(c config.Config, idp idempotent.Idempotent, log log.Logger) (mq.Subscriber, error) {
	address := make([]string, 0)
	err := c.Value("kafka.address").Scan(&address)
	if err != nil {
		return nil, err
	}
	topic, err := c.Value("kafka.topic").String()
	if err != nil {
		return nil, err
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  address,
		GroupID:  "group-a",
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	rlqPolicy := &RLQPolicy{
		RetryLetterTopic: defaultRetryTopic,
		Address:          address,
		GroupID:          "group-b",
	}
	retryRouter, err := newRetryRouter(rlqPolicy, r)
	if err != nil {
		return nil, err
	}

	dlqPolicy := &DLQPolicy{
		MaxDeliveries:   3,
		DeadLetterTopic: defaultDeadTopic,
		Address:         address,
		GroupID:         "group-b",
	}
	dlqRouter, err := newDlqRouter(dlqPolicy, r)
	if err != nil {
		return nil, err
	}
	k := &kafkaSubscriber{
		topic:  topic,
		reader: r,
		idp:    idp,
		rlq:    retryRouter,
		dlq:    dlqRouter,
	}
	go k.consumeRetryTopic(context.Background(), address)
	return k, nil
}

func (k *kafkaSubscriber) Subscribe(ctx context.Context) {
	for {
		m, err := k.reader.FetchMessage(context.Background())
		if err != nil {
			if err != io.EOF {
				logx.Error(ctx, "kafka fetch message fail", slog.Any("error", err))
			}
			continue
		}
		k.handleMessage(ctx, m, k.h)
	}
}

// handleMessage 幂等消费消息
func (k *kafkaSubscriber) handleMessage(ctx context.Context, m kafka.Message, h mq.Handler) {
	msg := NewMessage(&m)
	key := msg.PropsMessageID()
	ok, err := k.idp.SetKeyPendingWithDDL(ctx, key, time.Second*time.Duration(10*60))
	if err != nil {
		logx.Error(ctx, "subscriber SetKeyPendingWithDDL fail", slog.Any("error", err), slog.String("key", key), slog.Any("message", m))
		return
	}
	// 已经消费过
	if !ok {
		state, err := k.idp.GetKeyState(ctx, key)
		if err != nil {
			logx.Error(ctx, "subscriber GetKeyState fail", slog.Any("error", err), slog.String("key", key), slog.Any("message", m))
			return
		}
		if state == idempotent.Done {
			// 已经消费过，且消费成功了, 即使提交失败了，只要后面的message能成功，之前的message都会被提交
			if err = k.reader.CommitMessages(ctx, m); err != nil {
				logx.Error(ctx, "failed to commit messages:", slog.Any("error", err), slog.Any("message", m))
			}
		} else {
			// 已经消费过，但消费失败了
			k.ReconsumeLater(msg)
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
				return k.idp.RemoveFailKey(ctx, key)
			}, 2)
			if err != nil {
				logx.Error(ctx, "RemoveFailKey fail after retry 2 times", slog.Any("error", err), slog.String("key", key))
			}
			k.ReconsumeLater(msg)
		} else {
			// 消费采成功
			err = retry.Run(func() error {
				return k.idp.UpdateKeyDone(ctx, key)
			}, 2)
			if err != nil {
				logx.Error(ctx, "UpdateKeyDone fail after retry 2 times", slog.Any("error", err), slog.String("key", key))
			}
			err = retry.Run(func() error {
				return k.reader.CommitMessages(ctx, m)
			}, 2)
			if err != nil {
				logx.Error(ctx, "consumer ack fail after retry 2 times", slog.Any("error", err), slog.Any("message", m))
			}
		}
	}
}

func (k *kafkaSubscriber) Close() error {
	err := k.reader.Close()
	fmt.Print("close kafka reader")
	if err != nil {
		return err
	}
	return nil
}

// consumeRetryTopic 重新投递到原来的队列
func (k *kafkaSubscriber) consumeRetryTopic(ctx context.Context, address []string) {
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
			k.redelivery(ctx, r, w)
		}(topic)
	}
}

func (k *kafkaSubscriber) redelivery(ctx context.Context, r *kafka.Reader, w *kafka.Writer) {
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

func (k *kafkaSubscriber) ReconsumeLater(m *Message) {
	m.SetRealTopic(k.topic)
	m.IncrReconsumeTimes()
	m.SetDelayTime(retryBackoff(m.PropsReconsumeTimes()))
	if m.PropsReconsumeTimes() > k.dlq.maxRetry() {
		k.dlq.Chan() <- m
	} else {
		k.rlq.Chan() <- m
	}
}

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

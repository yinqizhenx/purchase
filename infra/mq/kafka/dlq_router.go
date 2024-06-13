package kafka

import (
	"context"
	"errors"
	"log/slog"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/segmentio/kafka-go"

	"purchase/infra/logx"
	"purchase/pkg/retry"
)

type dlqRouter struct {
	consumer  *kafka.Reader
	writer    *kafka.Writer
	policy    *DLQPolicy
	messageCh chan *Message
	closeCh   chan interface{}
	log       log.Logger
}

// DLQPolicy represents the configuration for the Dead Letter Queue consumer policy.
type DLQPolicy struct {
	// MaxDeliveries specifies the maximum number of times that a message will be delivered before being
	// sent to the dead letter queue.
	MaxDeliveries int

	// DeadLetterTopic specifies the name of the topic where the failing messages will be sent.
	DeadLetterTopic string

	GroupID string
	Address []string
}

func newDlqRouter(policy *DLQPolicy, logger log.Logger, rawReader *kafka.Reader) (*dlqRouter, error) {
	if policy == nil {
		return nil, errors.New("policy can not be nil")
	}
	r := &dlqRouter{
		policy: policy,
		log:    logger,
	}

	if policy != nil {
		if policy.MaxDeliveries <= 0 {
			return nil, errors.New("DLQPolicy.MaxDeliveries needs to be > 0")
		}

		if policy.DeadLetterTopic == "" {
			return nil, errors.New("DLQPolicy.Topic needs to be set to a valid topic name")
		}

		r.messageCh = make(chan *Message)
		r.closeCh = make(chan interface{}, 1)
		r.log = logger
		r.writer = &kafka.Writer{
			// Topic:    policy.DeadLetterTopic,
			Addr:                   kafka.TCP(policy.Address...),
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		}
		go r.run()
	}
	r.consumer = rawReader
	return r, nil
}

func (r *dlqRouter) shouldSendToDlq(cm *Message) bool {
	if r.policy == nil {
		return false
	}
	return cm.PropsReconsumeTimes() != 0 && cm.PropsReconsumeTimes() >= r.policy.MaxDeliveries
}

func (r *dlqRouter) maxRetry() int {
	return r.policy.MaxDeliveries
}

func (r *dlqRouter) Chan() chan *Message {
	return r.messageCh
}

func (r *dlqRouter) run() {
	for {
		select {
		case rm := <-r.messageCh:
			msg, err := rm.ToKafkaMessage(defaultDeadTopic)
			if err != nil {
				logx.Error(nil, "message transfer to kafka message fail", slog.Any("message", rm), slog.Any("error", err))
				break
			}
			err = retry.Run(func() error {
				return r.writer.WriteMessages(context.Background(), *msg)
			}, 3, retry.NewDefaultBackoffPolicy())

			if err != nil {
				logx.Error(nil, "message send to dead queue fail after retry 3 times", slog.Any("message", rm), slog.Any("error", err))
				break
			}

			msg.Topic = rm.PropsRealTopic()
			err = r.consumer.CommitMessages(context.Background(), *msg)
			if err != nil {
				logx.Error(nil, "commit message fail", slog.Any("message", rm), slog.Any("error", err))
			}
		case <-r.closeCh:
			if r.writer != nil {
				r.writer.Close()
			}
			return
		}
	}
}

func (r *dlqRouter) close() {
	// Attempt to write on the close channel, without blocking
	select {
	case r.closeCh <- nil:
	default:
	}
}

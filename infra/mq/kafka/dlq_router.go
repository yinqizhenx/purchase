package kafka

import (
	"context"
	"errors"
	"log/slog"

	"github.com/go-kratos/kratos/v2/log"

	"purchase/infra/logx"
	"purchase/infra/mq"
	"purchase/pkg/retry"
)

type dlqRouter struct {
	committer mq.MessageCommitter
	pub       mq.Publisher
	policy    *DLQPolicy
	messageCh chan *mq.Message
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

func newDlqRouter(address []string, committer mq.MessageCommitter, pub mq.Publisher) (*dlqRouter, error) {
	policy := &DLQPolicy{
		MaxDeliveries:   3,
		DeadLetterTopic: defaultDeadTopic,
		Address:         address,
		GroupID:         "group-b",
	}

	r := &dlqRouter{
		policy:    policy,
		committer: committer,
		pub:       pub,
	}

	if policy.MaxDeliveries <= 0 {
		return nil, errors.New("DLQPolicy.MaxDeliveries needs to be > 0")
	}

	if policy.DeadLetterTopic == "" {
		return nil, errors.New("DLQPolicy.Topic needs to be set to a valid topic name")
	}

	r.messageCh = make(chan *mq.Message)
	r.closeCh = make(chan interface{}, 1)
	// r.writer = &kafka.Writer{
	// 	// Topic:    policy.DeadLetterTopic,
	// 	Addr:                   kafka.TCP(policy.Address...),
	// 	Balancer:               &kafka.LeastBytes{},
	// 	AllowAutoTopicCreation: true,
	// }
	go r.run()
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

func (r *dlqRouter) Chan() chan *mq.Message {
	return r.messageCh
}

func (r *dlqRouter) run() {
	for {
		select {
		case msg := <-r.messageCh:
			// msg, err := rm.ToKafkaMessage(defaultDeadTopic)
			// if err != nil {
			// 	logx.Error(nil, "message transfer to kafka message fail", slog.Any("message", rm), slog.Any("error", err))
			// 	break
			// }
			err := retry.Run(func() error {
				return r.pub.Publish(context.Background(), msg)
			}, 3, retry.NewDefaultBackoffPolicy())

			if err != nil {
				logx.Error(nil, "message send to dead queue fail after retry 3 times", slog.Any("message", msg), slog.Any("error", err))
				break
			}

			// msg.Topic = rm.PropsRealTopic()
			err = r.committer.CommitMessage(context.Background(), msg)
			if err != nil {
				logx.Error(nil, "commit message fail", slog.Any("message", msg), slog.Any("error", err))
			}
		case <-r.closeCh:
			if r.pub != nil {
				r.pub.Close()
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

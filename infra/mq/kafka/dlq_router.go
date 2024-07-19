package kafka

import (
	"context"
	"errors"
	"log/slog"

	"github.com/go-kratos/kratos/v2/log"

	"purchase/infra/logx"
	"purchase/infra/mq"
)

type dlqRouter struct {
	// committer mq.MessageCommitter
	pub       mq.Publisher
	policy    *DLQPolicy
	messageCh chan RetryMessage
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

func newDlqRouter(address []string, pub mq.Publisher) (*dlqRouter, error) {
	policy := &DLQPolicy{
		MaxDeliveries:   3,
		DeadLetterTopic: defaultDeadTopic,
		Address:         address,
		GroupID:         "group-b",
	}

	r := &dlqRouter{
		policy: policy,
		// committer: committer,
		pub: pub,
	}

	if policy.MaxDeliveries <= 0 {
		return nil, errors.New("DLQPolicy.MaxDeliveries needs to be > 0")
	}

	if policy.DeadLetterTopic == "" {
		return nil, errors.New("DLQPolicy.Topic needs to be set to a valid topic name")
	}

	r.messageCh = make(chan RetryMessage)
	r.closeCh = make(chan interface{}, 1)
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

func (r *dlqRouter) Chan() chan RetryMessage {
	return r.messageCh
}

func (r *dlqRouter) run() {
	for {
		select {
		case msg := <-r.messageCh:
			msg.SetDeadTopic(defaultDeadTopic)

			err := r.pub.Publish(context.Background(), msg.Message)
			if err != nil {
				logx.Error(nil, "message send to dead queue fail after retry 3 times", slog.Any("message", msg), slog.Any("error", err))
				break
			}

			err = msg.Commit(context.Background())
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

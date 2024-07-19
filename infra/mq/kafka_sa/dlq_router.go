package kafka_sa

import (
	"context"
	"log/slog"

	"purchase/infra/logx"
	"purchase/infra/mq"
)

type dlqRouter struct {
	pub       mq.Publisher
	messageCh chan RetryMessage
	closeCh   chan interface{}
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

func newDlqRouter(pub mq.Publisher) *dlqRouter {

	r := &dlqRouter{
		pub: pub,
	}

	r.messageCh = make(chan RetryMessage)
	r.closeCh = make(chan interface{}, 1)
	go r.run()
	return r
}

func (r *dlqRouter) maxRetry() int {
	return defaultMaxRetry
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

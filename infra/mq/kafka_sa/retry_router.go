package kafka_sa

import (
	"context"
	"log/slog"
	"time"

	"purchase/infra/logx"
	"purchase/infra/mq"
)

type retryRouter struct {
	pub mq.Publisher
	// policy    *RLQPolicy
	messageCh chan RetryMessage
	closeCh   chan interface{}
}

type RLQPolicy struct {
	GroupID          string
	Address          []string
	RetryLetterTopic string
}

func newRetryRouter(pub mq.Publisher) *retryRouter {
	r := &retryRouter{
		pub: pub,
	}

	r.messageCh = make(chan RetryMessage)
	r.closeCh = make(chan interface{}, 1)

	go r.run()

	return r
}

func (r *retryRouter) Chan() chan RetryMessage {
	return r.messageCh
}

func (r *retryRouter) run() {
	for {
		select {
		case msg := <-r.messageCh:
			rtyTopic := genRetryTopic(msg.DelayTime())
			msg.SetRetryTopic(rtyTopic)

			err := r.pub.Publish(context.Background(), msg.Message)
			if err != nil {
				logx.Error(nil, "message send to retry queue fail", slog.Any("message", msg), slog.Any("error", err))
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
			// r.log.Debug("Closed RLQ router")
			return
		}
	}
}

func (r *retryRouter) close() {
	// Attempt to write on the close channel, without blocking
	select {
	case r.closeCh <- nil:
	default:
	}
}

func genRetryTopic(d time.Duration) string {
	return retryTopic[d]
}

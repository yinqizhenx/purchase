package kafka

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"purchase/infra/logx"
	"purchase/infra/mq"
)

type MessageCommitter interface {
	CommitMessages(ctx context.Context, msgs ...Message) error
}

type retryRouter struct {
	// client    kafka.Dialer
	committer mq.MessageCommitter
	pub       mq.Publisher
	policy    *RLQPolicy
	messageCh chan *mq.Message
	closeCh   chan interface{}
}

type RLQPolicy struct {
	GroupID string
	Address []string
	// RetryLetterTopic specifies the name of the topic where the retry messages will be sent.
	RetryLetterTopic string
}

func newRetryRouter(address []string, committer mq.MessageCommitter, pub mq.Publisher) (*retryRouter, error) {
	policy := &RLQPolicy{
		RetryLetterTopic: defaultRetryTopic,
		Address:          address,
		GroupID:          "group-b",
	}

	r := &retryRouter{
		policy:    policy,
		committer: committer,
		pub:       pub,
	}

	if policy.RetryLetterTopic == "" {
		return nil, errors.New("DLQPolicy.RetryLetterTopic needs to be set to a valid topic name")
	}
	r.messageCh = make(chan *mq.Message)
	r.closeCh = make(chan interface{}, 1)
	// r.log = logger
	// r.writer = &kafka.Writer{
	// 	Addr:                   kafka.TCP(policy.Address...),
	// 	Balancer:               &kafka.LeastBytes{},
	// 	AllowAutoTopicCreation: true,
	// }
	go r.run()

	return r, nil
}

func (r *retryRouter) Chan() chan *mq.Message {
	return r.messageCh
}

func (r *retryRouter) run() {
	for {
		select {
		case msg := <-r.messageCh:
			rtyTopic := genRetryTopic(msg.DelayTime())
			msg.SetRetryTopic(rtyTopic)

			// msg, err := rm.ToKafkaMessage(rtyTopic)
			// if err != nil {
			// 	logx.Error(nil, "message transfer to kafka message fail", slog.Any("message", rm), slog.Any("error", err))
			// 	break
			// }

			err := r.pub.Publish(context.Background(), msg)
			if err != nil {
				logx.Error(nil, "message send to retry queue fail", slog.Any("message", msg), slog.Any("error", err))
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

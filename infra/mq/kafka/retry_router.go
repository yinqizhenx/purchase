package kafka

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	// "github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/segmentio/kafka-go"

	"purchase/infra/logx"
)

type retryRouter struct {
	// client    kafka.Dialer
	consumer  *kafka.Reader
	writer    *kafka.Writer
	policy    *RLQPolicy
	messageCh chan *Message
	closeCh   chan interface{}
	log       log.Logger
}

type RLQPolicy struct {
	GroupID string
	Address []string
	// RetryLetterTopic specifies the name of the topic where the retry messages will be sent.
	RetryLetterTopic string
}

func newRetryRouter(policy *RLQPolicy, logger log.Logger, rawReader *kafka.Reader) (*retryRouter, error) {
	if policy == nil {
		return nil, errors.New("policy can not be nil")
	}
	r := &retryRouter{
		policy: policy,
		log:    logger,
	}

	if policy.RetryLetterTopic == "" {
		return nil, errors.New("DLQPolicy.RetryLetterTopic needs to be set to a valid topic name")
	}
	r.consumer = rawReader
	r.messageCh = make(chan *Message)
	r.closeCh = make(chan interface{}, 1)
	r.log = logger
	r.writer = &kafka.Writer{
		Addr:                   kafka.TCP(policy.Address...),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
	go r.run()

	return r, nil
}

func (r *retryRouter) Chan() chan *Message {
	return r.messageCh
}

func (r *retryRouter) run() {
	for {
		select {
		case rm := <-r.messageCh:
			rtyTopic := genRetryTopic(rm.PropsDelayTime())
			rm.SetRetryTopic(rtyTopic)
			
			msg, err := rm.ToKafkaMessage(rtyTopic)
			if err != nil {
				logx.Error(nil, "message transfer to kafka message fail", slog.Any("message", rm), slog.Any("error", err))
				break
			}

			err = r.writer.WriteMessages(context.Background(), *msg)
			if err != nil {
				logx.Error(nil, "message send to retry queue fail", slog.Any("message", rm), slog.Any("error", err))
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

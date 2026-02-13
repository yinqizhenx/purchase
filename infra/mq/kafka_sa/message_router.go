package kafka_sa

import (
	"context"
	"log/slog"

	"purchase/infra/logx"
	"purchase/infra/mq"
	"purchase/pkg/retry"
)

// topicSetter 设置消息要发送到的目标 topic
type topicSetter func(msg *RetryMessage)

type messageRouter struct {
	pub        mq.Publisher
	messageCh  chan RetryMessage
	closeCh    chan struct{}
	name       string
	setTopic   topicSetter
	maxRetries int
}

func newMessageRouter(name string, pub mq.Publisher, setTopic topicSetter) *messageRouter {
	r := &messageRouter{
		pub:        pub,
		messageCh:  make(chan RetryMessage),
		closeCh:    make(chan struct{}, 1),
		name:       name,
		setTopic:   setTopic,
		maxRetries: 3,
	}
	go r.run()
	return r
}

func (r *messageRouter) Chan() chan RetryMessage {
	return r.messageCh
}

func (r *messageRouter) run() {
	for {
		select {
		case msg := <-r.messageCh:
			r.setTopic(&msg)
			logx.Info(context.Background(), r.name+" received message",
				slog.String("bizCode", msg.BizCode()),
				slog.String("eventName", msg.EventName()),
			)
			err := retry.Run(func() error {
				return r.pub.Publish(context.Background(), msg.Message)
			}, r.maxRetries)
			if err != nil {
				// 发布失败，记录完整消息信息以便人工介入恢复
				// 此时不 commit 原始消息的 offset，Kafka rebalance 后消息会被重新消费
				logx.Error(context.Background(), r.name+" publish fail after retries, message may need manual recovery",
					slog.Any("error", err),
					slog.String("bizCode", msg.BizCode()),
					slog.String("eventName", msg.EventName()),
					slog.String("messageID", msg.ID),
					slog.Int("reconsumeTime", msg.ReconsumeTimes()),
				)
				break
			}
			err = msg.Commit(context.Background())
			if err != nil {
				logx.Error(context.Background(), r.name+" commit fail",
					slog.Any("error", err),
					slog.String("bizCode", msg.BizCode()),
				)
			}
		case <-r.closeCh:
			return
		}
	}
}

func (r *messageRouter) close() {
	select {
	case r.closeCh <- struct{}{}:
	default:
	}
}

func (r *messageRouter) maxRetry() int {
	return defaultMaxRetry
}

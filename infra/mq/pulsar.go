package mq

import (
	"context"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/go-kratos/kratos/v2/log"

	"purchase/infra/idempotent"
	"purchase/pkg/retry"
)

var _ Publisher = (*pulsarPublisher)(nil)
var _ Subscriber = (*pulsarSubscriber)(nil)

type pulsarPublisher struct {
	client   pulsar.Client
	producer pulsar.Producer
	idg      IDGenFunc
}

func NewPulsarPublisher(idg IDGenFunc) (Publisher, error) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://localhost:6650",
	})
	if err != nil {
		return nil, err
	}

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: "my-topic",
	})
	pub := &pulsarPublisher{
		client:   client,
		producer: producer,
		idg:      idg,
	}
	return pub, nil
}

func (pub *pulsarPublisher) Publish(ctx context.Context, m *Message) error {
	// val, err := e.Encode()
	// if err != nil {
	// 	return err
	// }
	msg := &pulsar.ProducerMessage{
		// Key:   []byte(e.EntityID()),
		Payload:    m.Body,
		Properties: m.Header(),
	}
	_, err := pub.producer.Send(context.Background(), msg)
	return err
}

func (pub *pulsarPublisher) Close() error {
	pub.client.Close()
	return nil
}

type pulsarSubscriber struct {
	client   pulsar.Client
	consumer pulsar.Consumer
	idp      idempotent.Idempotent
	sem      chan struct{} // 限制最大并发消费
}

func NewPulsarSubscriber(idp idempotent.Idempotent) (Subscriber, error) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:              "pulsar://localhost:6650",
		OperationTimeout: 30 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	d := &pulsar.DLQPolicy{
		MaxDeliveries:    3,
		RetryLetterTopic: "persistent://group/server/xxx-RETRY",
		DeadLetterTopic:  "persistent://group/server/xxx-DLQ",
	}
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:               "persistent://group/server/xxx",
		SubscriptionName:    "test",
		Type:                pulsar.Shared,
		RetryEnable:         true,
		DLQ:                 d,
		NackRedeliveryDelay: time.Second * 3,
	})
	sub := &pulsarSubscriber{
		client:   client,
		consumer: consumer,
		sem:      make(chan struct{}, 5),
		idp:      idp,
	}
	return sub, nil
}

func (sub pulsarSubscriber) Subscribe(ctx context.Context, h Handler) {
	// todo 定时清除长时间未消费成功的消息幂等键
	for {
		msg, err := sub.consumer.Receive(ctx)
		if err != nil {
			return
		}
		go func() {
			sub.sem <- struct{}{}
			defer func() {
				if p := recover(); p != nil {
					log.Errorf("consumer panic when handle message: %s, err: %v", msg, p)
				}
				<-sub.sem
			}()
			key := msg.Properties()[MessageID]
			ok, err := sub.idp.SetKeyPendingWithDDL(ctx, key, time.Minute*time.Duration(10))
			if err != nil {
				log.Error(err)
				return
			}
			// 已经消费过
			if !ok {
				s, err := sub.idp.GetKeyState(ctx, key)
				if err != nil {
					log.Error(err)
					return
				}
				if s == idempotent.Done {
					// 已经消费过，且消费成功了
					err = sub.consumer.Ack(msg)
					if err != nil {
						log.Error(err)
						return
					}
				} else {
					// 已经消费过，但消费失败了，投入重试队列
					sub.consumer.ReconsumeLater(msg, time.Second*5)
				}
			} else {
				// 未消费过，执行消费逻辑
				m := &Message{
					Body: msg.Payload(),
				}
				m.HeaderSet(EventName, msg.Properties()[EventName])
				err = h(ctx, m)
				if err != nil {
					// 消费失败, 投入重试队列
					// 如果RemoveFailKey失败，ReconsumeLater成功，重会等到key过期删除，然后被消费
					// 如果RemoveFailKey成功，ReconsumeLater失败，队列未提交ack, 会再次消费
					// 如果RemoveFailKey失败，ReconsumeLater失败，队列未提交ack，会等到key过期删除，然后被消费
					err = retry.Run(func() error {
						return sub.idp.RemoveFailKey(ctx, key)
					}, 2)
					if err != nil {
						log.Errorf("RemoveFailKey fail after retry 2 times :%v", err)
					}
					sub.consumer.ReconsumeLater(msg, time.Second*5)
				} else {
					// 消费采成功，ack
					// 如果UpdateKeyDone失败，ack成功，正常流转
					// 如果UpdateKeyDone成功，ack失败，再次消费时会直接ack
					// 如果UpdateKeyDone失败，ack失败，会等到key过期删除，然后被重复消费，需要人工介入
					err = retry.Run(func() error {
						return sub.idp.UpdateKeyDone(ctx, key)
					}, 2)
					if err != nil {
						log.Errorf("UpdateKeyDone fail after retry 2 times :%v", err)
					}
					err = retry.Run(func() error {
						return sub.consumer.Ack(msg)
					}, 2)
					if err != nil {
						log.Errorf("consumer ack fail after retry 2 times :%v", err)
					}
				}
			}
		}()
		// if msg.Key() == 0 {
		// 	// 确认的处理
		// 	sub.sub.Ack(msg)
		// } else {
		// 	// 不确认，等 NackRedeliveryDelay 后将被重新投递到主队列进行消费
		// 	sub.sub.Nack(msg)
		//
		// 	// 稍后处理,等 xx 秒后将被重新投递到重试队列
		// 	sub.sub.ReconsumeLater(msg, time.Second*5)
		//
		// 	// 以上方法二选其一
		// }
	}
}

func (sub pulsarSubscriber) Close() error {
	sub.client.Close()
	return nil
}

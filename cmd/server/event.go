package server

import (
	"purchase/adapter/consumer"
	"purchase/infra/mq"
)

func NewEventConsumerServer(s mq.Subscriber) *consumer.EventConsumer {
	return consumer.NewEventConsumer(s)
}

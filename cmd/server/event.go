package server

import (
	"purchase/adapter/consumer"
	"purchase/app"
	"purchase/infra/mq"
)

func NewEventConsumerServer(s mq.Subscriber, srv *app.DomainEventAppService) *consumer.EventConsumer {
	return consumer.NewEventConsumer(s, srv)
}

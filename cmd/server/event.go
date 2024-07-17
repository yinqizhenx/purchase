package server

import (
	"purchase/adapter/consumer"
	"purchase/app/event_handler"
	"purchase/infra/mq"
)

func NewEventConsumerServer(s mq.Subscriber, srv *event_handler.DomainEventAppService) *consumer.EventConsumer {
	return consumer.NewEventConsumer(s, srv)
}

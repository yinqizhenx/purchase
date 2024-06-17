package server

import (
	"purchase/adapter/consumer/event"
	"purchase/infra/mq"
)

func NewDomainEventServer(s mq.Subscriber) *event.DomainEventApp {
	return event.NewDomainEventApp(s)
}

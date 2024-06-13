package server

import (
	"purchase/app/event"
	"purchase/infra/mq"
)

func NewDomainEventServer(s mq.Subscriber) *event.DomainEventApp {
	return event.NewDomainEventApp(s)
}

package consumer

import (
	"context"

	"purchase/app"
	"purchase/infra/mq"
)

type EventConsumer struct {
	sub        mq.Subscriber
	appService *app.DomainEventAppService
}

func NewEventConsumer(s mq.Subscriber, srv *app.DomainEventAppService) *EventConsumer {
	ec := &EventConsumer{
		sub:        s,
		appService: srv,
	}
	return ec
}

func (s *EventConsumer) Start(ctx context.Context) error {
	s.sub.Subscribe(ctx, s.Consume)
	return nil
}

func (s *EventConsumer) Consume(ctx context.Context, m *mq.Message) error {
	return s.appService.Handle(ctx, m)
}

func (s *EventConsumer) Stop(ctx context.Context) error {
	return s.sub.Close()
}

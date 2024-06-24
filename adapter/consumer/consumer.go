package consumer

import (
	"context"

	"purchase/app"
	"purchase/infra/mq"
)

type EventConsumer struct {
	sub        mq.Subscriber
	appService *app.DomainEventAppService
	mdw        []MiddleWare
}

func NewEventConsumer(s mq.Subscriber, srv *app.DomainEventAppService) *EventConsumer {
	ec := &EventConsumer{
		sub:        s,
		appService: srv,
	}
	// ec.use()
	return ec
}

func (s *EventConsumer) Start(ctx context.Context) error {
	s.sub.Subscribe(ctx, s.Consume)
	return nil
}

func (s *EventConsumer) Consume(ctx context.Context, m *mq.Message) error {
	h := s.appService.Handle
	for i := len(s.mdw); i >= 0; i-- {
		h = s.mdw[i](h)
	}
	return h(ctx, m)
}

func (s *EventConsumer) use(m MiddleWare) {
	s.mdw = append(s.mdw, m)
}

func (s *EventConsumer) Stop(ctx context.Context) error {
	return s.sub.Close()
}

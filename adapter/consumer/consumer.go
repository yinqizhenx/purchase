package consumer

import (
	"context"

	"github.com/go-kratos/kratos/v2/config"

	"purchase/app/event_handler"
	"purchase/infra/mq"
)

type EventConsumer struct {
	// sub        mq.Subscriber
	appService *event_handler.DomainEventAppService
	// mdw        []MiddleWare
	// address    []string
}

func NewEventConsumer(s mq.Subscriber, srv *event_handler.DomainEventAppService, cfg config.Config) *EventConsumer {
	ec := &EventConsumer{
		// sub:        s,
		appService: srv,
	}
	// ec.use()
	return ec
}

func (s *EventConsumer) Start(ctx context.Context) error {
	// s.sub.Subscribe(ctx, s.Consume)
	s.appService.Handle(ctx)
	return nil
}

// func (s *EventConsumer) Consume(ctx context.Context, m *mq.Message) error {
// 	h := s.appService.Handle
// 	for i := len(s.mdw) - 1; i >= 0; i-- {
// 		h = s.mdw[i](h)
// 	}
// 	return h(ctx, m)
// }

// func (s *EventConsumer) use(m MiddleWare) {
// 	s.mdw = append(s.mdw, m)
// }

func (s *EventConsumer) Stop(ctx context.Context) error {
	// return s.sub.Close()
	return nil
}

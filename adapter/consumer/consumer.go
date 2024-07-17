package consumer

import (
	"context"

	"purchase/infra/mq/kafka"
)

type EventConsumer struct {
	sub *kafka.Sub
	// appService *event_handler.DomainEventHandler
	// mdw        []MiddleWare
	// address    []string
}

func NewEventConsumer(sub *kafka.Sub) *EventConsumer {
	ec := &EventConsumer{
		sub: sub,
		// appService: srv,
	}
	// ec.use()
	return ec
}

func (s *EventConsumer) Start(ctx context.Context) error {
	// s.sub.Subscribe(ctx, s.Consume)
	s.sub.Subscribe(ctx)
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

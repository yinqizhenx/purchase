package event_handler

import (
	"context"
	"fmt"

	domainEvent "purchase/domain/event"
	"purchase/infra/mq/kafka"
)

// type Handler func(context.Context, domainEvent.Event) error

type DomainEventAppService struct {
	// handlers map[domainEvent.Event][]domainEvent.Handler
	subscriber *kafka.Sub
}

// 在我看来，事件不过是一种特殊的 Command，与应用层作为外部请求的入口一样，事件的消费入口同样是在应用层
// DomainEventAppService 里的每个方法，都是对特定某个领域事件的处理。
// 方法的参数一般是 Context 和对应监听的领域事件，而返回值只是一个error，用来标识当前处理是否成功

func NewDomainEventAppService(sub *kafka.Sub) *DomainEventAppService {
	app := &DomainEventAppService{
		// handlers: make(map[domainEvent.Event][]domainEvent.Handler),
		subscriber: sub,
	}
	app.registerEventHandler(&domainEvent.PACreated{}, OnPACreatedHandler{})
	return app
}

func (s *DomainEventAppService) registerEventHandler(e domainEvent.Event, h domainEvent.Handler) {
	// s.handlers[e] = append(s.handlers[e], h)
	s.subscriber.RegisterEventHandler(e, h)
}

func (s *DomainEventAppService) Handle(ctx context.Context) error {
	// for e, hs := range s.handlers {
	// 	if e.EventName() == m.HeaderGet(mq.EventName) {
	// 		msg, err := e.Decode(m.Body)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		for _, h := range hs {
	// 			err = h.Handle(ctx, msg)
	// 			if err != nil {
	// 				logx.Errorf(ctx, "handle event fail; message: %v, err: %v", *m, err)
	// 			}
	// 		}
	// 		// return h(ctx, msg)
	// 	}
	// }
	// logx.Errorf(ctx, "no handler found for event [%s]", m.HeaderGet(mq.EventName))
	// return nil
	s.subscriber.Subscribe(ctx)
	return nil
}

func (s *DomainEventAppService) OnPACreated(ctx context.Context, e domainEvent.Event) error {
	if created, ok := e.(*domainEvent.PACreated); ok {
		fmt.Println("created event", created)
		// return errors.New("an error happened 11")
	}
	return nil
}

type OnPACreatedHandler struct{}

func (h OnPACreatedHandler) Handle(ctx context.Context, e domainEvent.Event) error {
	return nil
}

func (h OnPACreatedHandler) Name() string {
	return "pa_created_handler"
}

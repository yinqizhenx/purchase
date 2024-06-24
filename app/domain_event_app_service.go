package app

import (
	"context"
	"fmt"

	domainEvent "purchase/domain/event"
	"purchase/infra/logx"
	"purchase/infra/mq"
)

type Handler func(context.Context, domainEvent.Event) error

type DomainEventAppService struct {
	handlers map[domainEvent.Event]Handler
}

// 在我看来，事件不过是一种特殊的 Command，与应用层作为外部请求的入口一样，事件的消费入口同样是在应用层
// DomainEventAppService 里的每个方法，都是对特定某个领域事件的处理。
// 方法的参数一般是 Context 和对应监听的领域事件，而返回值只是一个error，用来标识当前处理是否成功

func NewDomainEventAppService() *DomainEventAppService {
	app := &DomainEventAppService{
		handlers: make(map[domainEvent.Event]Handler),
	}
	app.registerEventHandler(&domainEvent.PACreated{}, app.OnPACreated)
	return app
}

func (s *DomainEventAppService) registerEventHandler(e domainEvent.Event, h Handler) {
	s.handlers[e] = h
}

func (s *DomainEventAppService) Handle(ctx context.Context, m *mq.Message) error {
	for e, h := range s.handlers {
		if e.EventName() == m.HeaderGet(mq.EventName) {
			msg, err := e.Decode(m.Body)
			if err != nil {
				return err
			}
			return h(ctx, msg)
		}
	}
	logx.Errorf(ctx, "no handler found for event [%s]", m.HeaderGet(mq.EventName))
	return nil
}

func (s *DomainEventAppService) OnPACreated(ctx context.Context, e domainEvent.Event) error {
	if created, ok := e.(*domainEvent.PACreated); ok {
		fmt.Println("created event", created)
		// return errors.New("an error happened 11")
	}
	return nil
}

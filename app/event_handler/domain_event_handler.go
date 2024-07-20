package event_handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/wire"

	domainEvent "purchase/domain/event"
)

var ProviderSet = wire.NewSet(NewDomainEventHandler)

type DomainEventHandler struct {
	handlers map[domainEvent.Event][]domainEvent.Handler
}

// 在我看来，事件不过是一种特殊的 Command，与应用层作为外部请求的入口一样，事件的消费入口同样是在应用层
// DomainEventHandler 里的每个方法，都是对特定某个领域事件的处理。
// 方法的参数一般是 Context 和对应监听的领域事件，而返回值只是一个error，用来标识当前处理是否成功

func NewDomainEventHandler() domainEvent.HandlerAggregator {
	app := &DomainEventHandler{
		handlers: make(map[domainEvent.Event][]domainEvent.Handler),
	}
	app.registerEventHandler(&domainEvent.PACreated{}, app.OnPACreated)
	return app
}

func (s *DomainEventHandler) registerEventHandler(e domainEvent.Event, h domainEvent.Handler) {
	s.handlers[e] = append(s.handlers[e], h)
}

func (s *DomainEventHandler) Build() map[domainEvent.Event][]domainEvent.Handler {
	return s.handlers
}

func (s *DomainEventHandler) OnPACreated(ctx context.Context, e domainEvent.Event) error {
	if created, ok := e.(*domainEvent.PACreated); ok {
		fmt.Println("created event", created)
		return errors.New("an error happened 11")
	}
	return nil
}

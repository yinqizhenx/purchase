package event

import (
	"context"
	"time"

	"purchase/domain/entity/async_task"
)

// Event 包含一些通用的方法
type Event interface {
	ID() string                                  // 单个event的id
	EntityID() string                            // 领域事件的产生，一般是由于聚合状态的变更，因此，在领域事件上，还应该包含对应的聚合根id
	OccurredOn() time.Time                       // event 发生事件
	EventName() string                           // event name
	Encode() ([]byte, error)                     // 序列化event
	Decode(v []byte) (Event, error)              // 反序列化event
	ToAsyncTask() (*async_task.AsyncTask, error) // 将event转为异步任务
}

type PAEvent interface {
	Event
	OrderID() int64
}

type ProductEvent interface {
	Event
	ProductID() int64
}

type Handler interface {
	Handle(context.Context, Event) error
	Name() HandlerType
}

type HandlerType string

func (h HandlerType) String() string {
	return string(h)
}

type HandlerAggregator interface {
	Build() map[Event][]Handler
}

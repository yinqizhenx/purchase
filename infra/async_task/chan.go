package async_task

import (
	"context"

	"github.com/google/wire"

	"purchase/domain/entity/async_task"
	"purchase/pkg/chanx"
)

var ProviderSet = wire.NewSet(NewMessageChan)

const initCapacity = 10

func NewMessageChan() (*chanx.UnboundedChan[string], func()) {
	ctx, cancel := context.WithCancel(context.Background())
	return chanx.NewUnboundedChan[string](ctx, initCapacity), cancel
}

func NewTaskChan() (*chanx.UnboundedChan[*async_task.AsyncTask], func()) {
	ctx, cancel := context.WithCancel(context.Background())
	return chanx.NewUnboundedChan[*async_task.AsyncTask](ctx, initCapacity), cancel
}

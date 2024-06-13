package async_task

import (
	"context"

	"github.com/google/wire"

	"purchase/pkg/chanx"
)

var ProviderSet = wire.NewSet(NewMessageChan)

const initCapacity = 10

func NewMessageChan() (*chanx.UnboundedChan[string], func()) {
	ctx, cancel := context.WithCancel(context.Background())
	return chanx.NewUnboundedChan[string](ctx, initCapacity), cancel
}

package consumer

import (
	"context"

	"purchase/infra/mq"
)

type EventHandler func(ctx context.Context, m *mq.Message) error

type MiddleWare func(EventHandler) EventHandler

package repo

import (
	"context"

	"purchase/domain/event"
)

type EventRepo interface {
	Save(context.Context, ...event.Event) error // 保存一个event
}

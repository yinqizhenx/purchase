package idempotent

import (
	"context"
	"time"

	"github.com/google/wire"
)

// type State string

const (
	Done    string = "done"
	Pending string = "pending"
)

type Idempotent interface {
	SetKeyPendingWithDDL(ctx context.Context, key string, ddl time.Duration) (bool, error)
	GetKeyState(ctx context.Context, key string) (string, error)
	UpdateKeyDone(ctx context.Context, key string) error
	RemoveFailKey(ctx context.Context, key string) error
}

var ProviderSet = wire.NewSet(NewIdempotentImpl, NewMysqlIdempotentImpl)

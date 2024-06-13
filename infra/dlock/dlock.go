package dlock

import (
	"context"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

var ProviderSet = wire.NewSet(NewRedisLock)

type LockBuilder interface {
	NewLock(string) DistributeLock
}

type DistributeLock interface {
	Lock(context.Context) error
	Unlock(context.Context) error
}

var _ LockBuilder = (*RedisLock)(nil)

type RedisLock struct {
	rdb redis.UniversalClient
}

func NewRedisLock(rdb *redis.Client) LockBuilder {
	r := &RedisLock{
		rdb: rdb,
	}
	return r
}

func (r *RedisLock) NewLock(name string) DistributeLock {
	return NewMutex(r.rdb, name)
}

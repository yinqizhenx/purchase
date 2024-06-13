package idempotent

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ Idempotent = (*RedisIdempotentImpl)(nil)

type RedisIdempotentImpl struct {
	rdb *redis.Client
}

func NewIdempotentImpl(rdb *redis.Client) Idempotent {
	idp := &RedisIdempotentImpl{
		rdb: rdb,
	}
	return idp
}

func (imp RedisIdempotentImpl) SetKeyPendingWithDDL(ctx context.Context, key string, ddl time.Duration) (bool, error) {
	res := imp.rdb.SetNX(ctx, key, Pending, ddl)
	if err := res.Err(); err != nil {
		return false, err
	}
	return res.Val(), nil
}

func (imp RedisIdempotentImpl) GetKeyState(ctx context.Context, key string) (string, error) {
	s, err := imp.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return s, nil
}

func (imp RedisIdempotentImpl) UpdateKeyDone(ctx context.Context, key string) error {
	// 已完成的key保留3天
	res := imp.rdb.Set(ctx, key, Done, time.Hour*time.Duration(72))
	return res.Err()
}

func (imp RedisIdempotentImpl) RemoveFailKey(ctx context.Context, key string) error {
	res := imp.rdb.Del(ctx, key)
	return res.Err()
}

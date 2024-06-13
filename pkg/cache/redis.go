package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Codec interface {
	// Marshal returns the wire format of v.
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal parses the wire format into v.
	Unmarshal(data []byte, v interface{}) error
	// Name returns the name of the Codec implementation. The returned string
	// will be used as part of content type in transmission.  The result must be
	// static; the result cannot change between calls.
	Name() string
}

var _ Cache = (*redisCache)(nil)

type redisCache struct {
	rdb   *redis.Client
	codec Codec
}

func NewRedisCache(addr string) Cache {
	cli := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return redisCache{
		rdb: cli,
	}
}

func (c redisCache) Get(ctx context.Context, key string) (*Item, error) {
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	item := &Item{}
	err = c.codec.Unmarshal(data, item)
	return item, err
}

// Set stores a key-value pair into cache.
func (c redisCache) Set(ctx context.Context, key string, val interface{}, d time.Duration) error {
	item := Item{
		Value:      val,
		Expiration: time.Now().Add(d).UnixNano(),
	}
	data, err := c.codec.Marshal(item)
	if err != nil {
		return err
	}
	cmd := c.rdb.Set(ctx, key, data, d)
	return cmd.Err()
}

// Delete removes a key from cache.
func (c redisCache) Delete(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}

// String returns the name of the implementation.
func (c redisCache) String() string {
	return "redis"
}

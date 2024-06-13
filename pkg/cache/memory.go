package cache

import (
	"context"
	"time"
	"unsafe"

	"github.com/dgraph-io/ristretto"
)

// var Cache *memCache
//
// func InitCache() {
//	Cache = NewmemCache()
// }

type memCache struct {
	Cache *ristretto.Cache
	// urls  map[string][]string // url :[]schema
}

var _ Cache = &memCache{}

func NewMemCache() Cache {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
		Cost:        func(value interface{}) int64 { return int64(unsafe.Sizeof(value)) },
	})
	if err != nil {
		panic(err)
	}
	c := memCache{
		Cache: cache,
		// urls:  make(map[string][]string),
	}
	return &c
}

func (c *memCache) Set(ctx context.Context, key string, val interface{}, d time.Duration) error {
	var exp int64
	if d > 0 {
		exp = time.Now().Add(d).UnixNano()
	}
	item := Item{
		Value:      val,
		Expiration: exp,
	}
	// todo 处理返回false情况
	c.Cache.SetWithTTL(key, item, 0, d) // cost传0，则会使用Config.Cost动态计算占用内存大小
	c.Cache.Wait()
	return nil
}

func (c *memCache) Get(ctx context.Context, key string) (*Item, error) {
	val, ok := c.Cache.Get(key)
	if !ok {
		return nil, ErrKeyNotFound
	}
	item, ok := val.(Item)
	if !ok {
		return nil, ErrUnknownValue
	}
	return &item, nil
}

func (c *memCache) Delete(ctx context.Context, key string) error {
	c.Cache.Del(key)
	return nil
}

func (c *memCache) String() string {
	return "memory"
}

package db

import (
	"context"

	"purchase/infra/persistence/dal/db/ent"
	"purchase/pkg/cache"
)

// removeCache is a hook to remove related cache
func removeCache(cache cache.Cache) ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			value, err := next.Mutate(ctx, m)
			if err == nil { // todo 这里的错误会不会有不影响缓存更新的，
				// if m.Op() >= 1 && m.Op() <= 5 {
				// remove cache here
				// cache.DelBySchema(m.Type()) // todo 批量上传会导致缓存一直更新, 高并发时如何控制
				// }
			}
			return value, err
		})
	}
}

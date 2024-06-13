package dlock

import (
	"time"

	"github.com/go-redsync/redsync/v4"
)

type RedisDLock struct {
	lock *redsync.Redsync
}

func (rl *RedisDLock) Lock(s string) error {
	mutex := rl.newMutex(s)
	// 定时续期
	go func() {
		for range time.Tick(5 * time.Second) {
			if ok, err := mutex.Extend(); !ok || err != nil {
				return
			}
		}
	}()
	return mutex.Lock()
}

func (rl *RedisDLock) Unlock(s string) error {
	mutex := rl.newMutex(s)
	_, err := mutex.Unlock()
	return err
}

func (rl *RedisDLock) newMutex(s string) *redsync.Mutex {
	mutex := rl.lock.NewMutex(s, redsync.WithTries(1))
	return mutex
}

// func NewDistributeLock(rd *redis.Client) DistributeLock {
//
// 	pool := goredis.NewPool(rd) // or, pool := redigo.NewPool(...)
//
// 	// Create an instance of redisync to be used to obtain a mutual exclusion
// 	// lock.
// 	rs := redsync.New(pool)
// 	rl := &RedisDLock{
// 		lock: rs,
// 	}
// 	return rl
// }

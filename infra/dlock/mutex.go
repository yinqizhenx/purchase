package dlock

import (
	"context"
	"encoding/base64"
	"errors"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"

	"purchase/pkg/retry"
)

var ErrLockFailed = errors.New("redis lock: failed to acquire lock")

var ErrUnLockFailed = errors.New("redis lock: failed to release lock")

var ErrExtendExpiryFailed = errors.New("redis lock: failed to extend expiry")

var _ DistributeLock = (*Mutex)(nil)

type GenValueFunc func() (string, error)

const (
	defaultExpiry          = time.Duration(10 * time.Second)
	defaultTimeout         = time.Duration(500 * time.Millisecond)
	defaultTryCount        = 3
	defaultMaxExecDuration = time.Duration(10 * time.Minute)
)

type Mutex struct {
	name   string
	rdb    redis.UniversalClient
	expiry time.Duration
	// extExpInterval time.Duration // 续期定时器
	tryCount        int
	retry           retry.Retry
	timeout         time.Duration
	maxExecDuration time.Duration
	genValueFunc    GenValueFunc
	value           string
}

func NewMutex(rdb redis.UniversalClient, name string, opts ...Option) *Mutex {
	m := &Mutex{
		name:            name,
		rdb:             rdb,
		expiry:          defaultExpiry,
		tryCount:        defaultTryCount,
		retry:           retry.Run,
		timeout:         defaultTimeout,
		maxExecDuration: defaultMaxExecDuration,
		genValueFunc:    genValue,
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

func (m *Mutex) Lock(ctx context.Context) error {
	value, err := m.genValueFunc()
	if err != nil {
		return err
	}
	m.value = value
	err = m.retry(func() error {
		ctx, cancel := context.WithTimeout(ctx, m.timeout)
		defer cancel()
		ok, err := m.acquire(ctx, value)
		if err != nil {
			return err
		}
		if !ok {
			return ErrLockFailed
		}
		return nil
	}, m.tryCount)
	if err == nil {
		// 定时续期
		go func() {
			startTime := time.Now()
			for range time.Tick((m.expiry / 4) * 3) {
				// 超过最大执行时间，锁不再续期，避免锁删除失败，续期导致不退出
				if time.Since(startTime) > m.maxExecDuration {
					return
				}
				if err := m.ExtendExpiry(ctx); err != nil {
					return
				}
			}
		}()
	}
	return err
}

func (m *Mutex) Unlock(ctx context.Context) error {
	ok, err := m.release(ctx, m.value)
	if err != nil {
		return err
	}
	if !ok {
		return ErrUnLockFailed
	}
	return nil
}

func (m *Mutex) ExtendExpiry(ctx context.Context) error {
	ok, err := m.extendExpiry(ctx, m.value, int(m.expiry/time.Millisecond))
	if err != nil {
		return err
	}
	if !ok {
		return ErrExtendExpiryFailed
	}
	return nil
}

func (m *Mutex) acquire(ctx context.Context, value string) (bool, error) {
	reply, err := m.rdb.SetNX(ctx, m.name, value, m.expiry).Result()
	if err != nil {
		return false, err
	}
	return reply, nil
}

var deleteScript = `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`

func (m *Mutex) release(ctx context.Context, value string) (bool, error) {
	status, err := m.rdb.Eval(ctx, deleteScript, []string{m.name}, value).Result()
	if noErrNil(err) != nil {
		return false, err
	}
	return status != int64(0), nil
}

var extendExpiryScript = `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("PEXPIRE", KEYS[1], ARGV[2])
	else
		return 0
	end
`

func (m *Mutex) extendExpiry(ctx context.Context, value string, expiry int) (bool, error) {
	status, err := m.rdb.Eval(ctx, extendExpiryScript, []string{m.name}, value, expiry).Result()
	if noErrNil(err) != nil {
		return false, err
	}
	return status != int64(0), nil
}

func noErrNil(err error) error {
	if !errors.Is(err, redis.Nil) {
		return err
	}
	return errors.New("redis key not exist")
}

func genValue() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

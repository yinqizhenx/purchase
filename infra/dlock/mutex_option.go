package dlock

import (
	"time"

	"purchase/pkg/retry"
)

type Option func(mutex *Mutex)

func WithTryCount(n int) Option {
	return func(m *Mutex) {
		m.tryCount = n
	}
}

func WithRetryFunc(fn retry.Retry) Option {
	return func(m *Mutex) {
		m.retry = fn
	}
}

func WithExpiry(t time.Duration) Option {
	return func(m *Mutex) {
		m.expiry = t
	}
}

func WithTimeout(t time.Duration) Option {
	return func(m *Mutex) {
		m.timeout = t
	}
}

func WithGenValueFunc(fn GenValueFunc) Option {
	return func(m *Mutex) {
		m.genValueFunc = fn
	}
}

package retry

import (
	"math"
	"math/rand"
	"time"
)

type BackoffPolicy interface {
	Next(retryCount int) time.Duration
}

const (
	defaultInterval = time.Millisecond * time.Duration(200) // 200 ms
	defaultMaxDelay = time.Minute * time.Duration(10)       // 10 min
	defaultJitter   = 0.2
)

// defaultBackoffPolicy is default impl for BackoffPolicy.
type defaultBackoffPolicy struct {
	interval time.Duration
	maxDelay time.Duration
	// jitter provides a range to randomize backoff delays.
	jitter float64
}

func NewDefaultBackoffPolicy(opts ...Options) BackoffPolicy {
	backOff := &defaultBackoffPolicy{
		interval: defaultInterval,
		maxDelay: defaultMaxDelay,
		jitter:   defaultJitter,
	}
	for _, opt := range opts {
		opt(backOff)
	}
	if backOff.maxDelay <= backOff.interval {
		panic("maxDelay must larger than interval")
	}
	return backOff
}

func (bp *defaultBackoffPolicy) Next(retryCount int) time.Duration {
	if retryCount < 0 {
		return bp.interval
	}
	// 上下20%浮动
	dur := math.Abs(float64(bp.interval<<retryCount)) * bp.jitter * (rand.Float64()*2 - 1)
	return time.Duration(math.Min(dur, float64(bp.maxDelay)))
}

type Options func(*defaultBackoffPolicy)

func WithInterval(d time.Duration) Options {
	return func(o *defaultBackoffPolicy) {
		o.interval = d
	}
}

func WithMaxDelay(d time.Duration) Options {
	return func(o *defaultBackoffPolicy) {
		o.maxDelay = d
	}
}

func WithJitter(j float64) Options {
	return func(o *defaultBackoffPolicy) {
		if j <= 0 || j >= 1 {
			return
		}
		o.jitter = j
	}
}

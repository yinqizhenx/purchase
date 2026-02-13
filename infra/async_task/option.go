package async_task

import "time"

type Option func(*AsyncTaskMux)

func WithConcurrency(n int) Option {
	return func(m *AsyncTaskMux) {
		m.concurrency = n
	}
}

func WithMaxTaskLoad(n int) Option {
	return func(m *AsyncTaskMux) {
		m.maxTaskLoad = n
	}
}

func WithMiddleware(mws ...Middleware) Option {
	return func(m *AsyncTaskMux) {
		m.mdw = append(m.mdw, mws...)
	}
}

// WithShutdownTimeout 设置优雅关闭的超时时间，默认30秒
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(m *AsyncTaskMux) {
		m.shutdownTimeout = timeout
	}
}

// WithMaxRetry 设置任务失败自动重试最大次数，默认3次
func WithMaxRetry(n int) Option {
	return func(m *AsyncTaskMux) {
		m.maxRetry = n
	}
}

// WithStuckThreshold 设置任务卡住告警阈值，默认5分钟
func WithStuckThreshold(d time.Duration) Option {
	return func(m *AsyncTaskMux) {
		m.stuckThreshold = d
	}
}

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

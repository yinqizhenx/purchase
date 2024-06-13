package async_task

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

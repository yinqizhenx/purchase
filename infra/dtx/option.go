package dtx

import "time"

type Option func(*TransSaga)

func WithTimeout(d time.Duration) Option {
	return func(t *TransSaga) {
		t.timeout = d
	}
}

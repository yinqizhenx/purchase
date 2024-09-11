package dtx

import "time"

type Option func(*Saga)

func WithTimeout(d time.Duration) Option {
	return func(s *Saga) {
		s.timeout = d
	}
}

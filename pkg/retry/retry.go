package retry

import (
	"time"
)

type Retry func(fn func() error, n int, backoff ...BackoffPolicy) error

// Run n = 3 代表最多执行3次
func Run(fn func() error, n int, backoff ...BackoffPolicy) error {
	if n <= 0 {
		return fn()
	}
	bp := NewDefaultBackoffPolicy()
	if len(backoff) > 0 {
		bp = backoff[0]
	}
	var err error
	for i := 1; i <= n; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		time.Sleep(bp.Next(i))
	}
	return err
}

// func RunWithCallback(fn func() error, n int, onFail func(error) error, onSuccess func() error, backoff ...BackoffPolicy) error {
// 	err := Run(fn, n, backoff...)
// 	if err != nil {
// 		if onFail != nil {
// 			return onFail(err)
// 		}
// 		return err
// 	}
// 	if onSuccess != nil {
// 		return onSuccess()
// 	}
// 	return nil
// }

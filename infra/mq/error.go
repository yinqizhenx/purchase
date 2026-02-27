package mq

import "errors"

// RetryableError marks an error as worth retrying via the retry queue.
// By default, consumer errors go straight to the dead-letter topic.
// Only errors explicitly wrapped with NewRetryableError trigger ReconsumeLater.
type RetryableError struct {
	Err error
}

func (e *RetryableError) Error() string { return e.Err.Error() }
func (e *RetryableError) Unwrap() error { return e.Err }

func NewRetryableError(err error) error {
	return &RetryableError{Err: err}
}

func IsRetryable(err error) bool {
	var re *RetryableError
	return errors.As(err, &re)
}

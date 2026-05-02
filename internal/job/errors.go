package job

import "errors"

type RetryableError struct {
	cause error
}

func (e *RetryableError) Error() string { return e.cause.Error() }
func (e *RetryableError) Unwrap() error { return e.cause }

type PermanentError struct {
	cause error
}

func (e *PermanentError) Error() string { return e.cause.Error() }
func (e *PermanentError) Unwrap() error { return e.cause }

func NewRetryable(err error) error {
	return &RetryableError{cause: err}
}

func NewPermanent(err error) error {
	return &PermanentError{cause: err}
}

func IsRetryable(err error) bool {
	var e *RetryableError
	return errors.As(err, &e)
}

func IsPermanent(err error) bool {
	var e *PermanentError
	return errors.As(err, &e)
}

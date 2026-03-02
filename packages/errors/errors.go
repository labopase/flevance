package errors

import (
	"github.com/cockroachdb/errors"
)

type Error interface {
	error
	Code() int
	Message() string
}

type errx struct {
	error
	code    int
	message string
}

func (e *errx) Code() int {
	return e.code
}

func (e *errx) Message() string {
	return e.message
}

func (e *errx) Error() string {
	return e.error.Error()
}

func New(message string) error {
	return errors.New(message)
}

func Newf(format string, args ...interface{}) error {
	return errors.Newf(format, args...)
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

func WithMessage(err error, message string) error {
	return errors.WithMessage(err, message)
}

func WithMessagef(err error, format string, args ...interface{}) error {
	return errors.WithMessagef(err, format, args...)
}

func WithCode(code int, message string) error {
	return &errx{
		code:  code,
		error: New(message),
	}
}

func WithCodef(code int, format string, args ...interface{}) error {
	return &errx{
		code:  code,
		error: Newf(format, args...),
	}
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}

func GetCode(err error) int {
	if err == nil {
		return 0
	}

	var appErr Error
	if errors.As(err, &appErr) {
		return appErr.Code()
	}

	return 0
}

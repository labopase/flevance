package errors

import "net/http"

func NewHttpError(code int, format string, args ...interface{}) error {
	return WithCodef(code, format, args...)
}

func InternalServerError(message string) error {
	return WithCode(http.StatusInternalServerError, message)
}

func BadRequest(message string) error {
	return WithCode(http.StatusBadRequest, message)
}

func NotFound(message string) error {
	return WithCode(http.StatusNotFound, message)
}

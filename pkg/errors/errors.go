package errs

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int
	Message string
	Details any
}

type NotFoundError struct {
	*AppError
}

func (e *AppError) Error() string {
	if e.Details == nil {
		return fmt.Sprintf("%d: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("%d: %s (%#v)", e.Code, e.Message, e.Details)
}

func New(code int, msg string, items ...any) *AppError {
	var ctx any

	switch {
	case len(items) == 1:
		ctx = items[0]
	case len(items) > 1:
		ctx = items
	}

	return &AppError{
		Code:    code,
		Message: msg,
		Details: ctx,
	}
}

func Internal(msg string, items ...any) error {
	return New(http.StatusInternalServerError, msg, items...)
}

func NotFound(msg string, items ...any) error {
	return &NotFoundError{
		AppError: New(http.StatusNotFound, msg, items...),
	}
}

func Conflict(msg string, items ...any) error {
	return New(http.StatusConflict, msg, items...)
}

func BadRequest(msg string, items ...any) error {
	return New(http.StatusBadRequest, msg, items...)
}

func Unauthorized(msg string, items ...any) error {
	return New(http.StatusUnauthorized, msg, items...)
}

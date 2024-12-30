package apperrors

import (
	"fmt"
	"net/http"
)

var (
	ErrorUserNotFound       = NotFound("user not found")
	ErrorUserDataRequired   = BadRequest("user data required")
	ErrorUserFailedToCreate = Internal("failed to create user")
	ErrorUserFailedToGet    = Internal("failed to get user")
	ErrorUserFailedToUpdate = Internal("failed to update user")
	ErrorUserFailedToDelete = Internal("failed to delete user")

	ErrorSessionNotFound       = NotFound("session not found")
	ErrorSessionDataRequired   = BadRequest("session data is required")
	ErrorSessionFailedToCreate = Internal("failed to create session")
	ErrorSessionFailedToGet    = Internal("failed to get session")
	ErrorSessionFailedToUpdate = Internal("failed to update session")
	ErrorSessionFailedToDelete = Internal("failed to delete session")

	ErrorIncomeNotFound       = NotFound("income not found")
	ErrorIncomeDataRequired   = BadRequest("income data is required")
	ErrorIncomeFailedToCreate = Internal("failed to create income")
	ErrorIncomeFailedToGet    = Internal("failed to get income")
	ErrorIncomeFailedToUpdate = Internal("failed to update income")
	ErrorIncomeFailedToDelete = Internal("failed to delete income")

	ErrorTxFailedToBegin          = Internal("failed to begin transaction")
	ErrorUniqueConstraintViolated = Conflict("unique constraint violation")

	ErrorFailedToFetchUserId = BadRequest("failed to fetch user id")
)

type Error struct {
	Code    int
	Message string
	Details any
}

func (e *Error) Error() string {
	if e.Details == nil {
		return fmt.Sprintf("%d: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("%d: %s (%#v)", e.Code, e.Message, e.Details)
}

func New(code int, msg string, items ...any) *Error {
	var ctx any

	switch {
	case len(items) == 1:
		ctx = items[0]
	case len(items) > 1:
		ctx = items
	}

	return &Error{
		Code:    code,
		Message: msg,
		Details: ctx,
	}
}

func Internal(msg string, items ...any) error {
	return New(http.StatusInternalServerError, msg, items...)
}

func NotFound(msg string, items ...any) error {
	return New(http.StatusNotFound, msg, items...)
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

package api

import (
	"context"
	"errors"
	"github.com/ogen-go/ogen/ogenerrors"
	"incomster/backend/api/oas"
	"incomster/backend/api/validation"
	"incomster/backend/service"
	"incomster/config"
	errs "incomster/pkg/errors"
	"log"
	"net/http"
)

var (
	_ oas.Handler         = (*Handler)(nil)
	_ oas.SecurityHandler = (*Handler)(nil)
)

var (
	FailedToFetchUserId = errs.BadRequest("failed to fetch user id")
)

type Handler struct {
	Config    *config.Config
	Service   *service.Service
	Validator *validation.Validator
}

func NewHandler(config *config.Config, service *service.Service, validator *validation.Validator) *Handler {
	return &Handler{
		Config:    config,
		Service:   service,
		Validator: validator,
	}
}

func (h *Handler) NewError(ctx context.Context, err error) *oas.ErrorStatusCode {
	var appErr *errs.AppError
	if errors.As(err, &appErr) {
		return h.handleAppError(appErr)
	}

	if errors.Is(err, ogenerrors.ErrSecurityRequirementIsNotSatisfied) {
		return &oas.ErrorStatusCode{
			StatusCode: http.StatusUnauthorized,
			Response: oas.Error{
				Code:    http.StatusUnauthorized,
				Message: "security requirement is not satisfied",
			},
		}
	}

	msg := "internal server error"
	if h.Config.Env.IsDev() {
		msg = err.Error()
	}

	return &oas.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: oas.Error{
			Code:    http.StatusInternalServerError,
			Message: msg,
		},
	}
}

func (h *Handler) handleAppError(err *errs.AppError) *oas.ErrorStatusCode {
	if err.Details != nil {
		log.Print(err.Details)
	}

	return &oas.ErrorStatusCode{
		StatusCode: err.Code,
		Response: oas.Error{
			Code:    err.Code,
			Message: err.Message,
		},
	}
}

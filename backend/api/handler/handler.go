package api

import (
	"context"
	"errors"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/rs/zerolog/log"
	"incomster/backend/api/oas"
	"incomster/backend/api/validation"
	"incomster/backend/service"
	"incomster/config"
	"incomster/pkg/apperrors"
	"incomster/pkg/ctxutil"
	"incomster/pkg/ternary"
	"net/http"
)

var _ oas.Handler = (*Handler)(nil)

type Handler struct {
	AccountHandler
	SelfHandler
	IncomeHandler
	SecurityHandler

	Config    *config.Config
	Service   *service.Service
	Validator *validation.Validator
}

func NewHandler(config *config.Config, service *service.Service, validator *validation.Validator) *Handler {
	return &Handler{
		AccountHandler:  *NewAccountHandler(service.Account, validator.Account),
		SelfHandler:     *NewSelfHandler(service.User, validator.User),
		IncomeHandler:   *NewIncomeHandler(service.Income, validator.Income),
		SecurityHandler: *NewSecurityHandler(service.Security),

		Config:    config,
		Service:   service,
		Validator: validator,
	}
}

func (h *Handler) NewError(ctx context.Context, err error) *oas.ErrorStatusCode {
	var (
		userId, _ = ctxutil.GetUserId(ctx)
		appErr    *apperrors.Error
	)

	if errors.As(err, &appErr) {
		return h.handleAppError(userId, appErr)
	}

	if errors.Is(err, ogenerrors.ErrSecurityRequirementIsNotSatisfied) {
		return h.handleSecurityError(userId, err)
	}

	return h.handleInternalError(userId, err)
}

func (h *Handler) handleAppError(userId int, err *apperrors.Error) *oas.ErrorStatusCode {
	log.Warn().Str("error_kind", "app error").Int("user_id", userId).Err(err).Send()

	return &oas.ErrorStatusCode{
		StatusCode: err.Code,
		Response: oas.Error{
			Code:    err.Code,
			Message: err.Message,
		},
	}
}

func (h *Handler) handleSecurityError(userId int, err error) *oas.ErrorStatusCode {
	log.Warn().Str("error_kind", "security error").Int("user_id", userId).Err(err).Send()

	return &oas.ErrorStatusCode{
		StatusCode: http.StatusUnauthorized,
		Response: oas.Error{
			Code:    http.StatusUnauthorized,
			Message: "security requirement is not satisfied",
		},
	}
}

func (h *Handler) handleInternalError(userId int, err error) *oas.ErrorStatusCode {
	log.Warn().Str("error_kind", "internal error").Int("user_id", userId).Err(err).Send()

	return &oas.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: oas.Error{
			Code:    http.StatusInternalServerError,
			Message: ternary.Func(h.Config.Env.IsDev(), err.Error(), "internal server error"),
		},
	}
}

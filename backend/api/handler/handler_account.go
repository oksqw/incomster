package api

import (
	"golang.org/x/net/context"
	"incomster/backend/api/oas"
	"incomster/backend/api/validation"
	"incomster/backend/dto/sessiondto"
	"incomster/backend/dto/userdto"
	"incomster/backend/service"
	"incomster/pkg/apperrors"
	"incomster/pkg/ctxutil"
)

var _ oas.AccountHandler = (*AccountHandler)(nil)

type AccountHandler struct {
	service   *service.AccountService
	validator *validation.AccountValidator
}

func NewAccountHandler(service *service.AccountService, validator *validation.AccountValidator) *AccountHandler {
	return &AccountHandler{service: service, validator: validator}
}

func (h *AccountHandler) Register(ctx context.Context, req *oas.UserRegisterRequest) (*oas.Session, error) {
	if err := h.validator.Register(req); err != nil {
		return nil, err
	}

	session, err := h.service.Register(ctx, userdto.RegisterToInput(req))
	if err != nil {
		return nil, err
	}

	return sessiondto.CoreToOas(session), nil
}

func (h *AccountHandler) Login(ctx context.Context, req *oas.UserLoginRequest) (*oas.Session, error) {
	if err := h.validator.Login(req); err != nil {
		return nil, err
	}

	session, err := h.service.Login(ctx, userdto.LoginToInput(req))
	if err != nil {
		return nil, err
	}

	return sessiondto.CoreToOas(session), nil
}

func (h *AccountHandler) Logout(ctx context.Context) error {
	userId, err := ctxutil.GetUserId(ctx)
	if err != nil {
		return apperrors.ErrorFailedToFetchUserId
	}

	return h.service.Logout(ctx, userId)
}

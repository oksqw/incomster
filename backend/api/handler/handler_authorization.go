package api

import (
	"golang.org/x/net/context"
	"incomster/backend/api/oas"
	"incomster/backend/dto/sessiondto"
	"incomster/backend/dto/userdto"
	"incomster/pkg/ctxutil"
)

func (h *Handler) Register(ctx context.Context, req *oas.UserRegisterRequest) (*oas.Session, error) {
	if err := h.Validator.User.Register(req); err != nil {
		return nil, err
	}

	session, err := h.Service.Account.Register(ctx, userdto.RegisterToInput(req))
	if err != nil {
		return nil, err
	}

	return sessiondto.CoreToOas(session), nil
}

func (h *Handler) Login(ctx context.Context, req *oas.UserLoginRequest) (*oas.Session, error) {
	if err := h.Validator.User.Login(req); err != nil {
		return nil, err
	}

	session, err := h.Service.Account.Login(ctx, userdto.LoginToInput(req))
	if err != nil {
		return nil, err
	}

	return sessiondto.CoreToOas(session), nil
}

func (h *Handler) Logout(ctx context.Context) error {
	userId, err := ctxutil.GetUserId(ctx)
	if err != nil {
		return FailedToFetchUserId
	}

	return h.Service.Account.Logout(ctx, userId)
}

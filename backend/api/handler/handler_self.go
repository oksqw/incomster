package api

import (
	"context"
	"incomster/backend/api/oas"
	"incomster/backend/api/validation"
	"incomster/backend/dto/userdto"
	"incomster/backend/service"
	"incomster/core"
	"incomster/pkg/apperrors"
	"incomster/pkg/ctxutil"
)

var _ oas.SelfHandler = (*SelfHandler)(nil)

type SelfHandler struct {
	service   *service.UserService
	validator *validation.UserValidator
}

func NewSelfHandler(service *service.UserService, validator *validation.UserValidator) *SelfHandler {
	return &SelfHandler{service: service, validator: validator}
}

func (h *SelfHandler) UpdateSelf(ctx context.Context, req *oas.UserUpdateRequest) (*oas.User, error) {
	err := h.validator.Update(req)
	if err != nil {
		return nil, err
	}

	userId, err := ctxutil.GetUserId(ctx)
	if err != nil {
		return nil, apperrors.ErrorFailedToFetchUserId
	}

	user, err := h.service.Update(ctx, userdto.UpdateToInput(req, userId))
	if err != nil {
		return nil, err
	}

	return userdto.CoreToOas(user), nil
}

func (h *SelfHandler) GetSelf(ctx context.Context) (*oas.User, error) {
	userId, err := ctxutil.GetUserId(ctx)
	if err != nil {
		return nil, apperrors.ErrorFailedToFetchUserId
	}

	user, err := h.service.Get(ctx, &core.UserGetInput{Id: &userId})
	if err != nil {
		return nil, err
	}

	return userdto.CoreToOas(user), nil
}

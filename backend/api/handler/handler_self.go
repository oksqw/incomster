package api

import (
	"context"
	"incomster/backend/api/oas"
	"incomster/backend/dto/userdto"
	"incomster/core"
	"incomster/pkg/ctxutil"
)

func (h *Handler) UpdateSelf(ctx context.Context, req *oas.UserUpdateRequest) (*oas.User, error) {
	err := h.Validator.User.Update(req)
	if err != nil {
		return nil, err
	}

	userId, err := ctxutil.GetUserId(ctx)
	if err != nil {
		return nil, FailedToFetchUserId
	}

	user, err := h.Service.User.Update(ctx, userdto.UpdateToInput(req, userId))
	if err != nil {
		return nil, err
	}

	return userdto.CoreToOas(user), nil
}

func (h *Handler) GetSelf(ctx context.Context) (*oas.User, error) {
	userId, err := ctxutil.GetUserId(ctx)
	if err != nil {
		return nil, FailedToFetchUserId
	}

	user, err := h.Service.User.Get(ctx, &core.UserGetInput{Id: &userId})
	if err != nil {
		return nil, err
	}

	return userdto.CoreToOas(user), nil
}

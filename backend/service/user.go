package service

import (
	"golang.org/x/net/context"
	"incomster/backend/store"
	"incomster/core"
	"incomster/pkg/jwt"
)

type UserService struct {
	store     store.IUserStore
	tokenizer *jwt.Tokenizer
}

func NewUserService(store store.IUserStore) *UserService {
	return &UserService{store: store}
}

func (u *UserService) Update(ctx context.Context, input *core.UserUpdateInput) (*core.User, error) {
	return u.store.Update(ctx, input)
}

func (u *UserService) Get(ctx context.Context, input *core.UserGetInput) (*core.User, error) {
	return u.store.Get(ctx, input)
}

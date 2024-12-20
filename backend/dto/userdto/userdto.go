package userdto

import (
	"incomster/backend/api/oas"
	"incomster/backend/store/postgres/dal"
	"incomster/core"
)

func CoreToOas(in *core.User) *oas.User {
	return &oas.User{
		ID:        in.ID,
		Username:  in.Username,
		Name:      in.Name,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
}

func DalToCore(user *dal.User) *core.User {
	return &core.User{
		ID:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func CreatToDal(in *core.UserCreateInput) *dal.User {
	return &dal.User{
		Username: in.Username,
		Password: in.Password,
		Name:     in.Name,
	}
}

func RegisterToInput(in *oas.UserRegisterRequest) *core.UserCreateInput {
	return &core.UserCreateInput{
		Username: in.Username,
		Password: in.Password,
		Name:     in.Name,
	}
}

func LoginToInput(in *oas.UserLoginRequest) *core.UserLoginInput {
	return &core.UserLoginInput{
		Username: in.Username,
		Password: in.Password,
	}
}

func UpdateToInput(in *oas.UserUpdateRequest, id int) *core.UserUpdateInput {
	out := &core.UserUpdateInput{Id: id}

	if in.Username.IsSet() {
		out.Username = &in.Username.Value
	}

	if in.Password.IsSet() {
		out.Password = &in.Password.Value
	}

	if in.Name.IsSet() {
		out.Name = &in.Name.Value
	}

	return out
}

package core

import (
	"time"
)

type User struct {
	ID        int       `json:"id"         required:"true"`
	Username  string    `json:"username"   required:"true"`
	Password  string    `json:"password"   required:"true"`
	Name      string    `json:"name"       required:"true"`
	Role      string    `json:"role"       required:"true"`
	CreatedAt time.Time `json:"created_at" required:"true"`
	UpdatedAt time.Time `json:"updated_at" required:"true"`
}

type UserCreateInput struct {
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
	Name     string `json:"name"     required:"true"`
}

type UserLoginInput struct {
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
}

type UserGetInput struct {
	Id       *int    `json:"id"`
	Username *string `json:"username"`
}

type UserUpdateInput struct {
	Id       int     `json:"id"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	Name     *string `json:"name"    `
}

type UserDeleteInput struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

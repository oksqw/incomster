package ctxutil

import (
	"fmt"
	"golang.org/x/net/context"
)

const (
	userIdKey   = "user_id"
	userRoleKey = "user_role"
)

var (
	ErrorValueNotFound      = fmt.Errorf("value not found")
	ErrorIncorrectValueType = fmt.Errorf("incorrect value type")
)

func WithUserId(ctx context.Context, userId int) context.Context {
	return context.WithValue(ctx, userIdKey, userId)
}

func WithUserRole(ctx context.Context, userRole string) context.Context {
	return context.WithValue(ctx, userRoleKey, userRole)
}

func GetUserId(ctx context.Context) (int, error) {
	value := ctx.Value(userIdKey)

	if value == nil {
		return 0, ErrorValueNotFound
	}

	id, ok := value.(int)
	if !ok {
		return 0, ErrorIncorrectValueType
	}

	return id, nil
}

func GetUserRole(ctx context.Context) (string, error) {
	value := ctx.Value(userRoleKey)
	if value == nil {
		return "", ErrorValueNotFound
	}

	role, ok := value.(string)
	if !ok {
		return "", ErrorIncorrectValueType
	}

	return role, nil
}

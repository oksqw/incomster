package store

import (
	"context"
	"incomster/backend/store/migrate"
	"incomster/core"
)

type IMigrator interface {
	Up(ctx context.Context) (from, to int, err error)
	Down(ctx context.Context, version int) (from, to int, err error)
	Migrations() ([]migrate.Migration, error)
}

type IStore interface {
	User() IUserStore
	Income() IIncomeStore
	Session() ISessionStore
}

type IUserStore interface {
	Create(ctx context.Context, user *core.UserCreateInput) (*core.User, error)
	Get(ctx context.Context, input *core.UserGetInput) (*core.User, error)
	Update(ctx context.Context, input *core.UserUpdateInput) (*core.User, error)
	Delete(ctx context.Context, input *core.UserDeleteInput) (*core.User, error)
}

type IIncomeStore interface {
	Create(ctx context.Context, income *core.IncomeCreateInput) (*core.Income, error)
	Get(ctx context.Context, id int) (*core.Income, error)
	Find(ctx context.Context, filter *core.IncomesFilter) (*core.Incomes, error)
	Update(ctx context.Context, income *core.IncomeUpdateInput) (*core.Income, error)
	Delete(ctx context.Context, id int) (*core.Income, error)
}

type ISessionStore interface {
	Create(ctx context.Context, input *core.SessionCreateInput) (*core.Session, error)
	Update(ctx context.Context, input *core.SessionUpdateInput) (*core.Session, error)
	Get(ctx context.Context, input *core.SessionGetInput) (*core.Session, error)
	Delete(ctx context.Context, input *core.SessionGetInput) (*core.Session, error)
}

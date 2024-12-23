package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"incomster/backend/dto/userdto"
	"incomster/backend/store"
	"incomster/backend/store/postgres/dal"
	"incomster/core"
	errs "incomster/pkg/errors"
)

var _ store.IUserStore = (*UserStore)(nil)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (u *UserStore) Create(ctx context.Context, input *core.UserCreateInput) (*core.User, error) {
	user := userdto.CreatToDal(input)
	err := user.Insert(ctx, u.db, boil.Infer())
	if err != nil {
		return nil, u.handleUserCreateError(err)
	}

	return userdto.DalToCore(user), nil
}

func (u *UserStore) Get(ctx context.Context, input *core.UserGetInput) (*core.User, error) {
	var mods []qm.QueryMod

	if input.Id != nil {
		mods = append(mods, qm.Where("id = ?", input.Id))
	}
	if input.Username != nil {
		mods = append(mods, qm.Where("username = ?", input.Username))
	}

	if len(mods) == 0 {
		return nil, ErrorUserDataRequired
	}

	found, err := dal.Users(mods...).One(ctx, u.db)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrorUserNotFound
	}
	if err != nil {
		return nil, ErrorUserFailedToGet
	}

	return userdto.DalToCore(found), nil
}

func (u *UserStore) Update(ctx context.Context, input *core.UserUpdateInput) (*core.User, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, ErrorTxFailedToBegin
	}
	defer CommitOrRollback(tx, err)

	found, err := dal.FindUser(ctx, tx, input.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrorUserNotFound
	}
	if err != nil {
		return nil, ErrorUserFailedToGet
	}

	var whitelist []string

	if input.Username != nil {
		found.Username = *input.Username
		whitelist = append(whitelist, dal.UserColumns.Username)
	}
	if input.Password != nil {
		found.Password = *input.Password
		whitelist = append(whitelist, dal.UserColumns.Password)
	}
	if input.Name != nil {
		found.Name = *input.Name
		whitelist = append(whitelist, dal.UserColumns.Name)
	}

	if len(whitelist) == 0 {
		return nil, ErrorUserDataRequired
	}

	_, err = found.Update(ctx, tx, boil.Whitelist(whitelist...))
	if err != nil {
		return nil, u.handleUserUpdateError(err)
	}

	return userdto.DalToCore(found), nil
}

func (u *UserStore) Delete(ctx context.Context, input *core.UserDeleteInput) (*core.User, error) {
	var mods []qm.QueryMod

	if input.Id > 0 {
		mods = append(mods, qm.Where("id = ?", input.Id))
	}
	if input.Username != "" {
		mods = append(mods, qm.Where("username = ?", input.Username))
	}
	if len(mods) == 0 {
		return nil, ErrorUserDataRequired
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, ErrorTxFailedToBegin
	}
	defer CommitOrRollback(tx, err)

	found, err := dal.Users(mods...).One(ctx, u.db)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrorUserNotFound
	}
	if err != nil {
		return nil, ErrorUserFailedToGet
	}

	_, err = found.Delete(ctx, tx)
	if err != nil {
		return nil, ErrorUserFailedToDelete
	}

	return userdto.DalToCore(found), nil
}

func (u *UserStore) isUniqueViolation(err error) (bool, string) {
	var pqerr *pq.Error
	ok := errors.As(err, &pqerr)

	if !ok {
		return false, ""
	}
	return pqerr.Code == "23505", pqerr.Constraint
}

func (u *UserStore) handleUserCreateError(err error) error {
	unique, constraint := u.isUniqueViolation(err)
	if !unique {
		return ErrorUserFailedToCreate
	}

	return u.handleUniqueConstraintsError(constraint)
}

func (u *UserStore) handleUserUpdateError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrorUserNotFound
	}

	unique, constraint := u.isUniqueViolation(err)
	if !unique {
		return ErrorUserFailedToUpdate
	}

	return u.handleUniqueConstraintsError(constraint)
}

func (u *UserStore) handleUniqueConstraintsError(constraint string) error {
	message, exists := uniqueConstraints[constraint]
	if exists {
		return errs.Conflict(message)
	}

	return ErrorUniqueConstraintViolated
}

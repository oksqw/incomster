package postgres

import (
	"incomster/backend/store"
	"incomster/pkg/errors"
)

var _ store.IStore = (*Store)(nil)

var (
	ErrorUserNotFound       = errs.NotFound("user not found")
	ErrorUserDataRequired   = errs.BadRequest("user data required")
	ErrorUserFailedToCreate = errs.Internal("failed to create user")
	ErrorUserFailedToGet    = errs.Internal("failed to get user")
	ErrorUserFailedToUpdate = errs.Internal("failed to update user")
	ErrorUserFailedToDelete = errs.Internal("failed to delete user")

	ErrorSessionNotFound       = errs.NotFound("session not found")
	ErrorSessionDataRequired   = errs.BadRequest("session data is required")
	ErrorSessionFailedToCreate = errs.Internal("failed to create session")
	ErrorSessionFailedToGet    = errs.Internal("failed to get session")
	ErrorSessionFailedToUpdate = errs.Internal("failed to update session")
	ErrorSessionFailedToDelete = errs.Internal("failed to delete session")

	ErrorIncomeNotFound       = errs.NotFound("income not found")
	ErrorIncomeDataRequired   = errs.BadRequest("income data is required")
	ErrorIncomeFailedToCreate = errs.Internal("failed to create income")
	ErrorIncomeFailedToGet    = errs.Internal("failed to get income")
	ErrorIncomeFailedToUpdate = errs.Internal("failed to update income")
	ErrorIncomeFailedToDelete = errs.Internal("failed to delete income")

	ErrorTxFailedToBegin  = errs.Internal("failed to begin transaction")
	ErrorTxFailedToCommit = errs.Internal("failed to commit transaction")

	ErrorUniqueConstraintViolated = errs.Conflict("unique constraint violation")

	uniqueConstraints = map[string]string{
		"users_username_key": "username is already taken",
	}
)

type Store struct {
	user    store.IUserStore
	income  store.IIncomeStore
	session store.ISessionStore
}

func NewStore(user store.IUserStore, income store.IIncomeStore, session store.ISessionStore) *Store {
	return &Store{user: user, income: income, session: session}
}

func (s *Store) User() store.IUserStore {
	return s.user
}

func (s *Store) Income() store.IIncomeStore {
	return s.income
}

func (s *Store) Session() store.ISessionStore { return s.session }

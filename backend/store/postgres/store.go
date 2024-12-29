package postgres

import (
	"incomster/backend/store"
)

var _ store.IStore = (*Store)(nil)

var uniqueConstraints = map[string]string{
	"users_username_key": "username is already taken",
}

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

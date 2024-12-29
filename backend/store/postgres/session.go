package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"incomster/backend/store"
	"incomster/backend/store/postgres/dal"
	"incomster/core"
	errs "incomster/pkg/apperrors"
)

var _ store.ISessionStore = (*SessionStore)(nil)

type SessionStore struct {
	db *sql.DB
}

func NewSessionStore(db *sql.DB) SessionStore {
	return SessionStore{db: db}
}

func (s SessionStore) Create(ctx context.Context, input *core.SessionCreateInput) (*core.Session, error) {
	// TODO: use ctx logger

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errs.ErrorTxFailedToBegin
	}
	defer CommitOrRollback(tx, err)

	session := &dal.Session{
		UserID: input.UserID,
		Token:  input.Token,
	}

	err = session.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return nil, errs.ErrorSessionFailedToCreate
	}

	mods := []qm.QueryMod{
		qm.Where("user_id = ?", session.UserID),
		qm.Where("token = ?", session.Token),
	}

	found, err := dal.Sessions(mods...).One(ctx, tx)
	if err != nil {
		return nil, errs.ErrorSessionFailedToGet
	}

	return s.toCore(found), nil
}

func (s SessionStore) Update(ctx context.Context, input *core.SessionUpdateInput) (*core.Session, error) {
	// TODO: use ctx logger

	session := &dal.Session{
		ID:     input.Id,
		UserID: input.UserID,
		Token:  input.Token,
	}

	_, err := session.Update(ctx, s.db, boil.Infer())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrorSessionNotFound
	}
	if err != nil {
		return nil, errs.ErrorSessionFailedToUpdate
	}

	return s.toCore(session), nil
}

func (s SessionStore) Get(ctx context.Context, input *core.SessionGetInput) (*core.Session, error) {
	// TODO: use ctx logger

	var mods []qm.QueryMod

	if input.Id > 0 {
		mods = append(mods, qm.Where("id = ?", input.Id))
	}
	if input.UserID > 0 {
		mods = append(mods, qm.Where("user_id = ?", input.UserID))
	}
	if input.Token != "" {
		mods = append(mods, qm.Where("token = ?", input.Token))
	}
	if len(mods) == 0 {
		return nil, errs.ErrorSessionDataRequired
	}

	found, err := dal.Sessions(mods...).One(ctx, s.db)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrorSessionNotFound
	}
	if err != nil {
		return nil, errs.ErrorSessionFailedToGet
	}

	return s.toCore(found), nil
}

func (s SessionStore) Delete(ctx context.Context, input *core.SessionGetInput) (*core.Session, error) {
	// TODO: use ctx logger

	var mods []qm.QueryMod

	if input.Id > 0 {
		mods = append(mods, qm.Where("id = ?", input.Id))
	}
	if input.UserID > 0 {
		mods = append(mods, qm.Where("user_id = ?", input.UserID))
	}
	if input.Token != "" {
		mods = append(mods, qm.Where("token = ?", input.Token))
	}
	if len(mods) == 0 {
		return nil, errs.ErrorSessionDataRequired
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errs.ErrorTxFailedToBegin
	}
	defer CommitOrRollback(tx, err)

	found, err := dal.Sessions(mods...).One(ctx, tx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrorSessionNotFound
	}
	if err != nil {
		return nil, errs.ErrorSessionFailedToGet
	}

	_, err = found.Delete(ctx, tx)
	if err != nil {
		return nil, errs.ErrorSessionFailedToDelete
	}

	return s.toCore(found), nil
}

func (s SessionStore) toCore(session *dal.Session) *core.Session {
	return &core.Session{
		ID:     session.ID,
		UserID: session.UserID,
		Token:  session.Token,
	}
}

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"incomster/backend/dto/incomedto"
	"incomster/backend/store"
	"incomster/backend/store/postgres/dal"
	"incomster/core"
	"incomster/pkg/apperrors"
	"incomster/pkg/collectionutils"
)

var _ store.IIncomeStore = (*IncomeStore)(nil)

type IncomeStore struct {
	db *sql.DB
}

func NewIncomeStore(db *sql.DB) *IncomeStore {
	return &IncomeStore{db: db}
}

func (i *IncomeStore) Create(ctx context.Context, input *core.IncomeCreateInput) (*core.Income, error) {
	income := incomedto.CreateToDal(input)
	err := income.Insert(ctx, i.db, boil.Infer())
	if err != nil {
		return nil, apperrors.ErrorIncomeFailedToCreate
	}

	return incomedto.DalToCore(income), nil
}

func (i *IncomeStore) Update(ctx context.Context, input *core.IncomeUpdateInput) (*core.Income, error) {
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, apperrors.ErrorTxFailedToBegin
	}
	defer CommitOrRollback(tx, err)

	found, err := dal.FindIncome(ctx, tx, input.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.ErrorIncomeNotFound
	}
	if err != nil {
		return nil, apperrors.ErrorIncomeFailedToGet
	}

	var whitelist []string

	if input.Amount != nil {
		found.Amount = *input.Amount
		whitelist = append(whitelist, dal.IncomeColumns.Amount)
	}
	if input.Comment != nil {
		found.Comment = null.StringFromPtr(input.Comment)
		whitelist = append(whitelist, dal.IncomeColumns.Comment)
	}

	if len(whitelist) == 0 {
		return nil, apperrors.ErrorIncomeDataRequired
	}

	_, err = found.Update(ctx, tx, boil.Whitelist(whitelist...))
	if err != nil {
		return nil, apperrors.ErrorIncomeFailedToUpdate
	}

	return incomedto.DalToCore(found), nil
}

func (i *IncomeStore) Delete(ctx context.Context, input *core.IncomeDeleteInput) (*core.Income, error) {
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, apperrors.ErrorTxFailedToBegin
	}
	defer CommitOrRollback(tx, err)

	mods := []qm.QueryMod{
		qm.Where(fmt.Sprintf("%s = ?", dal.IncomeColumns.ID), input.ID),
		qm.Where(fmt.Sprintf("%s = ?", dal.IncomeColumns.UserID), input.UserID),
	}

	found, err := dal.Incomes(mods...).One(ctx, i.db)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.ErrorIncomeNotFound
	}
	if err != nil {
		return nil, apperrors.ErrorIncomeFailedToGet
	}

	_, err = found.Delete(ctx, tx)
	if err != nil {
		return nil, apperrors.ErrorIncomeFailedToDelete
	}

	return incomedto.DalToCore(found), nil
}

func (i *IncomeStore) Get(ctx context.Context, input *core.IncomeGetInput) (*core.Income, error) {
	mods := []qm.QueryMod{
		qm.Where(fmt.Sprintf("%s = ?", dal.IncomeColumns.ID), input.ID),
		qm.Where(fmt.Sprintf("%s = ?", dal.IncomeColumns.UserID), input.UserID),
	}

	income, err := dal.Incomes(mods...).One(ctx, i.db)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.ErrorIncomeNotFound
	}
	if err != nil {
		return nil, apperrors.ErrorIncomeFailedToGet
	}

	return incomedto.DalToCore(income), nil
}

func (i *IncomeStore) Find(ctx context.Context, input *core.GetIncomesInput) (*core.Incomes, error) {
	mods := []qm.QueryMod{
		qm.Limit(input.Limit),
		qm.Offset(input.Offset),
	}

	if len(input.Users) > 0 {
		mods = append(mods, qm.WhereIn(fmt.Sprintf("%s IN ?", dal.IncomeColumns.UserID), collectionutils.ToInterface(input.Users)...))
	}
	if input.MinDate != nil {
		mods = append(mods, qm.Where(fmt.Sprintf("%s >= ?", dal.IncomeColumns.CreatedAt), input.MinDate))
	}
	if input.MaxDate != nil {
		mods = append(mods, qm.Where(fmt.Sprintf("%s <= ?", dal.IncomeColumns.CreatedAt), input.MaxDate))
	}
	if input.MinAmount != nil {
		mods = append(mods, qm.Where(fmt.Sprintf("%s >= ?", dal.IncomeColumns.Amount), input.MinAmount))
	}
	if input.MaxAmount != nil {
		mods = append(mods, qm.Where(fmt.Sprintf("%s <= ?", dal.IncomeColumns.Amount), input.MaxAmount))
	}

	found, err := dal.Incomes(mods...).All(ctx, i.db)

	if errors.Is(err, sql.ErrNoRows) {
		return &core.Incomes{}, nil
	}
	if err != nil {
		return nil, apperrors.Internal("failed to fetch incomes", err)
	}

	items := make([]*core.Income, len(found))
	count := len(found)
	sum := float64(0)

	for index, income := range found {
		items[index] = incomedto.DalToCore(income)
		sum += income.Amount
	}

	return &core.Incomes{
		Incomes: items,
		Count:   count,
		Sum:     sum,
	}, nil

}

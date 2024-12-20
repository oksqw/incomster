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
		return nil, ErrorIncomeFailedToCreate
	}

	return incomedto.DalToCore(income), nil
}

func (i *IncomeStore) Update(ctx context.Context, input *core.IncomeUpdateInput) (*core.Income, error) {
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, ErrorTxFailedToBegin
	}
	defer CommitOrRollback(tx, err)

	found, err := dal.FindIncome(ctx, tx, input.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrorIncomeNotFound
	}
	if err != nil {
		return nil, ErrorIncomeFailedToGet
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
		return nil, ErrorIncomeDataRequired
	}

	_, err = found.Update(ctx, tx, boil.Whitelist(whitelist...))
	if err != nil {
		return nil, ErrorIncomeFailedToUpdate
	}

	return incomedto.DalToCore(found), nil
}

func (i *IncomeStore) Delete(ctx context.Context, id int) (*core.Income, error) {
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, ErrorTxFailedToBegin
	}
	defer CommitOrRollback(tx, err)

	found, err := dal.FindIncome(ctx, tx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrorIncomeNotFound
	}
	if err != nil {
		return nil, ErrorIncomeFailedToGet
	}

	_, err = found.Delete(ctx, tx)
	if err != nil {
		return nil, ErrorIncomeFailedToDelete
	}

	return incomedto.DalToCore(found), nil
}

func (i *IncomeStore) Get(ctx context.Context, id int) (*core.Income, error) {
	income, err := dal.FindIncome(ctx, i.db, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrorIncomeNotFound
	}
	if err != nil {
		return nil, ErrorIncomeFailedToGet
	}

	return incomedto.DalToCore(income), nil
}

func (i *IncomeStore) Find(ctx context.Context, filter *core.IncomesFilter) (*core.Incomes, error) {
	mods := []qm.QueryMod{
		qm.Limit(filter.Limit),
		qm.Offset(filter.Offset),
	}

	if filter.MinDate != nil {
		mods = append(mods, qm.Where("created_at >= ?", filter.MinDate))
	}
	if filter.MaxDate != nil {
		mods = append(mods, qm.Where("created_at <= ?", filter.MaxDate))
	}
	if filter.MinAmount != nil {
		mods = append(mods, qm.Where("amount >= ?", filter.MinAmount))
	}
	if filter.MaxAmount != nil {
		mods = append(mods, qm.Where("amount <= ?", filter.MaxAmount))
	}

	found, err := dal.Incomes(mods...).All(ctx, i.db)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch incomes: %w", err)
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

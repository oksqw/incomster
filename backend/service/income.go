package service

import (
	"golang.org/x/net/context"
	"incomster/backend/store"
	"incomster/core"
)

type IncomeService struct {
	store store.IIncomeStore
}

func NewIncomeService(store store.IIncomeStore) *IncomeService {
	return &IncomeService{store: store}
}

func (s *IncomeService) Create(ctx context.Context, input *core.IncomeCreateInput) (*core.Income, error) {
	return s.store.Create(ctx, input)
}

func (s *IncomeService) Update(ctx context.Context, input *core.IncomeUpdateInput) (*core.Income, error) {
	return s.store.Update(ctx, input)
}

func (s *IncomeService) Delete(ctx context.Context, id int) (*core.Income, error) {
	return s.store.Delete(ctx, id)
}

func (s *IncomeService) Get(ctx context.Context, id int) (*core.Income, error) {
	return s.store.Get(ctx, id)
}

func (s *IncomeService) Find(ctx context.Context, filter *core.IncomesFilter) (*core.Incomes, error) {
	return s.store.Find(ctx, filter)
}

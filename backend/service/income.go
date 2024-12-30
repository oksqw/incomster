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

func (s *IncomeService) Delete(ctx context.Context, input *core.IncomeDeleteInput) (*core.Income, error) {
	return s.store.Delete(ctx, input)
}

func (s *IncomeService) Get(ctx context.Context, input *core.IncomeGetInput) (*core.Income, error) {
	return s.store.Get(ctx, input)
}

func (s *IncomeService) Find(ctx context.Context, input *core.GetIncomesInput) (*core.Incomes, error) {
	return s.store.Find(ctx, input)
}

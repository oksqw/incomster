package api

import (
	"context"
	"incomster/backend/api/oas"
	"incomster/backend/dto/incomedto"
	"incomster/backend/dto/incomesdto"
	"incomster/pkg/ctxutil"
	errs "incomster/pkg/errors"
)

func (h *Handler) AddIncome(ctx context.Context, req *oas.IncomeCreateRequest) (*oas.Income, error) {
	if err := h.Validator.Income.Create(req); err != nil {
		return nil, err
	}

	userId, err := ctxutil.GetUserId(ctx)
	if err != nil {
		return nil, errs.BadRequest("failed to get user id")
	}

	income, err := h.Service.Income.Create(ctx, incomedto.CreateToInput(req, userId))
	if err != nil {
		return nil, err
	}

	return incomedto.CoreToOas(income), nil
}

func (h *Handler) UpdateIncome(ctx context.Context, req *oas.IncomeUpdateRequest, params oas.UpdateIncomeParams) (*oas.Income, error) {
	if err := h.Validator.Income.Update(req); err != nil {
		return nil, err
	}

	income, err := h.Service.Income.Update(ctx, incomedto.UpdateToInput(req, params.ID))
	if err != nil {
		return nil, err
	}

	return incomedto.CoreToOas(income), nil
}

func (h *Handler) GetIncome(ctx context.Context, params oas.GetIncomeParams) (*oas.Income, error) {
	income, err := h.Service.Income.Get(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return incomedto.CoreToOas(income), nil
}

func (h *Handler) GetIncomes(ctx context.Context, params oas.GetIncomesParams) (*oas.Incomes, error) {
	incomes, err := h.Service.Income.Find(ctx, incomesdto.ParamsToFilter(&params))
	if err != nil {
		return nil, err
	}

	return incomesdto.CoreToOas(incomes), nil
}

package api

import (
	"context"
	"incomster/backend/api/oas"
	"incomster/backend/api/validation"
	"incomster/backend/dto/incomedto"
	"incomster/backend/dto/incomesdto"
	"incomster/backend/service"
	"incomster/pkg/apperrors"
	"incomster/pkg/ctxutil"
)

// TODO: Есть большая проблема.
// Методы взаимодействия с икомами не проверяют ID пользователя.
// Соответственно любой пользователь может взаимодействовать со всеми инкомами.
//----------------------------------------------------------------------------------------------------------------------
// Возможные варианты решения :
// Получать id пользователя из контекста и переделать некоторые структуры и методы стора
// чтобы в них учитывался не только по id инкома но и id пользователя.
// Таким образом при попытке получения чужого инкома пользователь будет получать ошибку 404.

var _ oas.IncomeHandler = (*IncomeHandler)(nil)

type IncomeHandler struct {
	service   *service.IncomeService
	validator *validation.IncomeValidator
}

func NewIncomeHandler(service *service.IncomeService, validator *validation.IncomeValidator) *IncomeHandler {
	return &IncomeHandler{service: service, validator: validator}
}

func (h *IncomeHandler) AddIncome(ctx context.Context, req *oas.IncomeCreateRequest) (*oas.Income, error) {
	if err := h.validator.Create(req); err != nil {
		return nil, err
	}

	userId, err := ctxutil.GetUserId(ctx)
	if err != nil {
		return nil, apperrors.ErrorFailedToFetchUserId
	}

	income, err := h.service.Create(ctx, incomedto.CreateToInput(req).WithUserId(userId))
	if err != nil {
		return nil, err
	}

	return incomedto.CoreToOas(income), nil
}

func (h *IncomeHandler) UpdateIncome(ctx context.Context, req *oas.IncomeUpdateRequest, params oas.UpdateIncomeParams) (*oas.Income, error) {
	if err := h.validator.Update(req); err != nil {
		return nil, err
	}

	userId, err := ctxutil.GetUserId(ctx)
	if err != nil {
		return nil, apperrors.ErrorFailedToFetchUserId
	}

	income, err := h.service.Update(ctx, incomedto.UpdateToInput(req).WithId(params.ID).WithUserId(userId))
	if err != nil {
		return nil, err
	}

	return incomedto.CoreToOas(income), nil
}

func (h *IncomeHandler) GetIncome(ctx context.Context, params oas.GetIncomeParams) (*oas.Income, error) {
	userId, err := ctxutil.GetUserId(ctx)
	if err != nil {
		return nil, apperrors.ErrorFailedToFetchUserId
	}

	income, err := h.service.Get(ctx, incomedto.GetParamsToInput(params).WithUserId(userId))
	if err != nil {
		return nil, err
	}

	return incomedto.CoreToOas(income), nil
}

func (h *IncomeHandler) GetIncomes(ctx context.Context, params oas.GetIncomesParams) (*oas.Incomes, error) {
	userId, err := ctxutil.GetUserId(ctx)
	if err != nil {
		return nil, apperrors.ErrorFailedToFetchUserId
	}

	incomes, err := h.service.Find(ctx, incomesdto.GetParamsToInput(&params).WithUserId(userId))
	if err != nil {
		return nil, err
	}

	return incomesdto.CoreToOas(incomes), nil
}

package incomedto

import (
	"github.com/volatiletech/null/v8"
	"incomster/backend/api/oas"
	"incomster/backend/store/postgres/dal"
	"incomster/core"
)

func CreateToInput(in *oas.IncomeCreateRequest) *core.IncomeCreateInput {
	out := &core.IncomeCreateInput{
		Amount: in.Amount,
	}

	if in.Comment.IsSet() {
		out.Comment = &in.Comment.Value
	}

	return out
}

func UpdateToInput(in *oas.IncomeUpdateRequest) *core.IncomeUpdateInput {
	out := &core.IncomeUpdateInput{}

	if in.Amount.IsSet() {
		out.Amount = &in.Amount.Value
	}

	if in.Comment.IsSet() {
		out.Comment = &in.Comment.Value
	}

	return out
}

func GetParamsToInput(in oas.GetIncomeParams) *core.IncomeGetInput {
	return &core.IncomeGetInput{
		ID: in.ID,
	}
}

func UpdateToDal(in *oas.IncomeUpdateRequest, incomeId int) *dal.Income {
	out := &dal.Income{ID: incomeId}

	if in.Amount.IsSet() {
		out.Amount = in.Amount.Value
	}

	if in.Comment.IsSet() {
		out.Comment = null.StringFrom(in.Comment.Value)
	}

	return out
}

func CreateToDal(in *core.IncomeCreateInput) *dal.Income {
	out := &dal.Income{
		UserID:  in.UserID,
		Amount:  in.Amount,
		Comment: null.StringFromPtr(in.Comment),
	}

	if in.CreatedAt != nil {
		out.CreatedAt = *in.CreatedAt
	}

	if in.UpdatedAt != nil {
		out.UpdatedAt = *in.UpdatedAt
	}

	return out
}

func CoreToOas(in *core.Income) *oas.Income {
	out := &oas.Income{
		ID:        in.ID,
		UserId:    in.UserID,
		Amount:    in.Amount,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}

	if in.Comment != "" {
		out.Comment.SetTo(in.Comment)
	}

	return out
}

func OasToCore(in *oas.Income) *core.Income {
	out := &core.Income{
		ID:        in.ID,
		UserID:    in.UserId,
		Amount:    in.Amount,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}

	if in.Comment.IsSet() {
		out.Comment = in.Comment.Value
	}

	return out
}

func DalToCore(in *dal.Income) *core.Income {
	return &core.Income{
		ID:        in.ID,
		UserID:    in.UserID,
		Amount:    in.Amount,
		Comment:   in.Comment.String,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
}

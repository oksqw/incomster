package incomesdto

import (
	"incomster/backend/api/oas"
	"incomster/backend/dto/incomedto"
	"incomster/core"
	"incomster/pkg/mapping"
)

func CoreToOas(in *core.Incomes) *oas.Incomes {
	return &oas.Incomes{
		Count:   in.Count,
		Sum:     in.Sum,
		Incomes: mapping.Map(in.Incomes, func(x *core.Income) oas.Income { return *incomedto.CoreToOas(x) }),
	}
}

func GetParamsToInput(in *oas.GetIncomesParams) *core.GetIncomesInput {
	out := &core.GetIncomesInput{
		Limit:  in.Limit,
		Offset: in.Offset,
	}

	if in.MinDate.IsSet() {
		out.MinDate = &in.MinDate.Value
	}

	if in.MaxDate.IsSet() {
		out.MaxDate = &in.MaxDate.Value
	}

	if in.MinAmount.IsSet() {
		out.MinAmount = &in.MinAmount.Value
	}

	if in.MaxAmount.IsSet() {
		out.MaxAmount = &in.MaxAmount.Value
	}

	return out
}

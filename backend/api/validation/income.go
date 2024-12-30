package validation

import (
	"incomster/backend/api/oas"
	errs "incomster/pkg/apperrors"
)

type IncomeValidator struct {
	Create func(*oas.IncomeCreateRequest) error
	Update func(*oas.IncomeUpdateRequest) error
}

func NewIncomeValidator() *IncomeValidator {
	return &IncomeValidator{
		Create: func(input *oas.IncomeCreateRequest) error {
			if input.Amount <= 0 {
				return errs.BadRequest("invalid amount")
			}
			if input.Comment.IsSet() && input.Comment.Value == "" {
				return errs.BadRequest("invalid comment")
			}
			return nil
		},
		Update: func(input *oas.IncomeUpdateRequest) error {
			if !input.Amount.IsSet() && !input.Comment.IsSet() {
				return errs.BadRequest("empty income data")
			}
			if input.Amount.IsSet() && input.Amount.Value <= 0 {
				return errs.BadRequest("invalid amount")
			}
			if input.Comment.IsSet() && input.Comment.Value == "" {
				return errs.BadRequest("invalid comment")
			}
			return nil
		},
	}
}

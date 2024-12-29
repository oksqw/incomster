package validation

import (
	"incomster/backend/api/oas"
	"incomster/pkg/apperrors"
	"incomster/pkg/passwordutil"
	"incomster/pkg/usernameutil"
)

type AccountValidator struct {
	Register func(*oas.UserRegisterRequest) error
	Login    func(*oas.UserLoginRequest) error
}

func NewAccountValidator() *AccountValidator {
	return &AccountValidator{
		Register: func(input *oas.UserRegisterRequest) error {
			if input.Username == "" {
				return apperrors.BadRequest("username is required")
			}
			if input.Name == "" {
				return apperrors.BadRequest("name is required")
			}
			if input.Password == "" {
				return apperrors.BadRequest("password is required")
			}
			if err := passwordutil.Validate(input.Password); err != nil {
				return apperrors.BadRequest(err.Error())
			}
			if err := usernameutil.Validate(input.Username); err != nil {
				return apperrors.BadRequest(err.Error())
			}
			return nil
		},
		Login: func(input *oas.UserLoginRequest) error {
			if input.Username == "" {
				return apperrors.BadRequest("username is required")
			}
			if input.Password == "" {
				return apperrors.BadRequest("password is required")
			}
			if len(input.Password) < 8 {
				return apperrors.BadRequest("password must be at least 8 characters")
			}
			return nil
		},
	}
}

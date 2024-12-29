package validation

import (
	"incomster/backend/api/oas"
	"incomster/pkg/errors"
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
				return errs.BadRequest("username is required")
			}
			if input.Name == "" {
				return errs.BadRequest("name is required")
			}
			if input.Password == "" {
				return errs.BadRequest("password is required")
			}
			if err := passwordutil.Validate(input.Password); err != nil {
				return errs.BadRequest(err.Error())
			}
			if err := usernameutil.Validate(input.Username); err != nil {
				return errs.BadRequest(err.Error())
			}
			return nil
		},
		Login: func(input *oas.UserLoginRequest) error {
			if input.Username == "" {
				return errs.BadRequest("username is required")
			}
			if input.Password == "" {
				return errs.BadRequest("password is required")
			}
			if len(input.Password) < 8 {
				return errs.BadRequest("password must be at least 8 characters")
			}
			return nil
		},
	}
}

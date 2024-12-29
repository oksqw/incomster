package validation

import (
	"incomster/backend/api/oas"
	"incomster/pkg/errors"
	"incomster/pkg/passwordutil"
	"incomster/pkg/usernameutil"
)

type UserValidator struct {
	Update func(*oas.UserUpdateRequest) error
}

func NewUserValidator() *UserValidator {
	return &UserValidator{
		Update: func(input *oas.UserUpdateRequest) error {
			if input.Username.IsSet() {
				if err := usernameutil.Validate(input.Username.Value); err != nil {
					return errs.BadRequest(err.Error())
				}
			}
			if input.Password.IsSet() {
				if err := passwordutil.Validate(input.Password.Value); err != nil {
					return errs.BadRequest(err.Error())
				}
			}
			return nil
		},
	}
}

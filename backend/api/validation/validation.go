package validation

type Validator struct {
	User    *UserValidator
	Income  *IncomeValidator
	Account *AccountValidator
}

func NewValidator() *Validator {
	return &Validator{
		User:    NewUserValidator(),
		Income:  NewIncomeValidator(),
		Account: NewAccountValidator(),
	}
}

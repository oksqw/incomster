package validation

type Validator struct {
	User   *UserValidator
	Income *IncomeValidator
}

func NewValidator() *Validator {
	return &Validator{
		User:   NewUserValidator(),
		Income: NewIncomeValidator(),
	}
}

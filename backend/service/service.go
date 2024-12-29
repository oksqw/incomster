package service

import errs "incomster/pkg/errors"

var (
	invalidCredential     = errs.Unauthorized("invalid username or password")
	failedToGenerateToken = errs.Internal("failed to generate JWT token")
	failedToValidateToken = errs.Internal("failed to validate JWT token")
	failedToUpdateToken   = errs.Internal("failed to update JWT token")
	failedToRetrieveUser  = errs.Internal("failed to retrieve user")
	failedToCreateSession = errs.Internal("failed to create session")
)

type Service struct {
	User     *UserService
	Income   *IncomeService
	Account  *AccountService
	Security *SecurityService
}

func NewService(user *UserService, income *IncomeService, account *AccountService, security *SecurityService) *Service {
	return &Service{User: user, Income: income, Account: account, Security: security}
}

package store

import errs "incomster/pkg/errors"

var (
	ErrorUserNotFound       = errs.NotFound("user not found")
	ErrorUserDataRequired   = errs.BadRequest("user data required")
	ErrorUserFailedToCreate = errs.Internal("failed to create user")
	ErrorUserFailedToGet    = errs.Internal("failed to get user")
	ErrorUserFailedToUpdate = errs.Internal("failed to update user")
	ErrorUserFailedToDelete = errs.Internal("failed to delete user")

	ErrorSessionNotFound       = errs.NotFound("session not found")
	ErrorSessionDataRequired   = errs.BadRequest("session data is required")
	ErrorSessionFailedToCreate = errs.Internal("failed to create session")
	ErrorSessionFailedToGet    = errs.Internal("failed to get session")
	ErrorSessionFailedToUpdate = errs.Internal("failed to update session")
	ErrorSessionFailedToDelete = errs.Internal("failed to delete session")

	ErrorIncomeNotFound       = errs.NotFound("income not found")
	ErrorIncomeDataRequired   = errs.BadRequest("income data is required")
	ErrorIncomeFailedToCreate = errs.Internal("failed to create income")
	ErrorIncomeFailedToGet    = errs.Internal("failed to get income")
	ErrorIncomeFailedToUpdate = errs.Internal("failed to update income")
	ErrorIncomeFailedToDelete = errs.Internal("failed to delete income")

	ErrorTxFailedToBegin  = errs.Internal("failed to begin transaction")
	ErrorTxFailedToCommit = errs.Internal("failed to commit transaction")

	ErrorUniqueConstraintViolated = errs.Conflict("unique constraint violation")
)

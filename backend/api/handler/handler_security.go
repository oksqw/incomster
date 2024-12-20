package api

import (
	"context"
	"incomster/backend/api/oas"
)

func (h *Handler) HandleBearerAuth(ctx context.Context, operationName oas.OperationName, t oas.BearerAuth) (context.Context, error) {
	return h.Service.Security.HandleBearerAuth(ctx, operationName, t.Token)
}

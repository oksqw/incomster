package api

import (
	"context"
	"incomster/backend/api/oas"
	"incomster/backend/service"
)

var _ oas.SecurityHandler = (*SecurityHandler)(nil)

type SecurityHandler struct {
	service *service.SecurityService
}

func NewSecurityHandler(service *service.SecurityService) *SecurityHandler {
	return &SecurityHandler{service: service}
}

func (h *SecurityHandler) HandleBearerAuth(ctx context.Context, operationName oas.OperationName, t oas.BearerAuth) (context.Context, error) {
	return h.service.HandleBearerAuth(ctx, operationName, t.Token)
}

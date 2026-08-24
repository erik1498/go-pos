package domain

import (
	"context"
	"go-pos/internal/model"
)

type AuditLogRepository interface {
	Create(ctx context.Context, logData model.AuditLog)
}

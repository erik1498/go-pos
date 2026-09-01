package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID        uuid.UUID
	ActorID   uuid.UUID
	ActorRole string
	Action    string
	Entity    string
	EntityID  string
	OldValues string
	NewValues string
	IPAddress string
	UserAgent string
	CreatedAt time.Time
}

type AuditLogRepository interface {
	Create(ctx context.Context, logData AuditLog) error
}

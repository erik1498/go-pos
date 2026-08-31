package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Tax struct {
	ID             uuid.UUID       `json:"id"`
	Name           string          `json:"name"`
	Rate           decimal.Decimal `json:"rate"`
	IsActive       bool            `json:"is_active"`
	IdempotencyKey string          `json:"idempotency_key"`
	CreatedBy      uuid.UUID       `json:"created_by"`
	UpdatedBy      uuid.UUID       `json:"updated_by"`
	DeletedBy      *uuid.UUID      `json:"deleted_by"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      *time.Time      `json:"deleted_at"`
}

type CreateTaxParam struct {
	Name string          `json:""`
	Rate decimal.Decimal `json:""`
}

type UpdateTaxParam struct {
	Name string          `json:""`
	Rate decimal.Decimal `json:""`
}

type TaxRepository interface {
	GetAll(ctx context.Context, opts QueryOptions) ([]Tax, int64, error)
	Create(ctx context.Context, tax Tax) (Tax, error)
	GetByID(ctx context.Context, id uuid.UUID) (Tax, error)
	UpdateByID(ctx context.Context, id uuid.UUID, tax Tax) (Tax, error)
	DeleteByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
}

type TaxUsecase interface {
	GetAll(ctx context.Context, opts QueryOptions) ([]Tax, int64, error)
	Create(ctx context.Context, req CreateTaxParam) (Tax, error)
	GetByID(ctx context.Context, id uuid.UUID) (Tax, error)
	UpdateByID(ctx context.Context, id uuid.UUID, tax UpdateTaxParam) (Tax, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

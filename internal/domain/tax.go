package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Tax struct {
	ID             uuid.UUID
	Name           string
	Rate           decimal.Decimal
	IsActive       bool
	IdempotencyKey string
	CreatedBy      uuid.UUID
	UpdatedBy      uuid.UUID
	DeletedBy      *uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

type CreateTaxParam struct {
	Name string
	Rate decimal.Decimal
}

type UpdateTaxParam struct {
	Name string
	Rate decimal.Decimal
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

package domain

import (
	"context"
	"errors"
	"go-pos/internal/model"

	"github.com/google/uuid"
)

var (
	ErrTaxNotFound = errors.New("tax data not found")
)

type TaxRepository interface {
	GetAll(ctx context.Context, opts QueryOptions) ([]model.Tax, int64, error)
	Create(ctx context.Context, tax model.Tax) (model.Tax, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Tax, error)
	UpdateByID(ctx context.Context, id uuid.UUID, tax model.Tax) (model.Tax, error)
	DeleteByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
}

type TaxUsecase interface {
	GetAll(ctx context.Context, opts QueryOptions) ([]model.Tax, int64, error)
	Create(ctx context.Context, req model.TaxRequest) (model.Tax, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Tax, error)
	UpdateByID(ctx context.Context, id uuid.UUID, tax model.TaxRequest) (model.Tax, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

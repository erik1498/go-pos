package domain

import (
	"context"
	"go-pos/internal/model"

	"github.com/google/uuid"
)

type ProductRepository interface {
	GetAll(ctx context.Context, opts QueryOptions) ([]model.Product, int64, error)
	Create(ctx context.Context, product model.Product) (model.Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Product, error)
	UpdateByID(ctx context.Context, id uuid.UUID, product model.Product) (model.Product, error)
	DeleteByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
}

type ProductUsecase interface {
	GetAll(ctx context.Context, opts QueryOptions) ([]model.Product, int64, error)
	Create(ctx context.Context, req model.ProductRequest) (model.Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Product, error)
	UpdateByID(ctx context.Context, id uuid.UUID, req model.ProductRequest) (model.Product, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

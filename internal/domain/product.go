package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID             uuid.UUID
	CategoryID     uuid.UUID
	Category       *Category
	Name           string
	SKU            string
	Price          decimal.Decimal
	Stock          decimal.Decimal
	Taxes          []Tax
	IdempotencyKey string
	CreatedBy      uuid.UUID
	UpdatedBy      uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

type CreateProductParam struct {
	CategoryID uuid.UUID
	Name       string
	SKU        string
	Price      decimal.Decimal
	Taxes      []Tax
}

type UpdateProductParam struct {
	CategoryID uuid.UUID
	Name       string
	SKU        string
	Price      decimal.Decimal
	Taxes      []Tax
}

type ProductRepository interface {
	GetAll(ctx context.Context, opts QueryOptions) ([]Product, int64, error)
	Create(ctx context.Context, product Product) (Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (Product, error)
	UpdateByID(ctx context.Context, id uuid.UUID, product Product) (Product, error)
	DeleteByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
}

type ProductUsecase interface {
	GetAll(ctx context.Context, opts QueryOptions) ([]Product, int64, error)
	Create(ctx context.Context, req CreateProductParam) (Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (Product, error)
	UpdateByID(ctx context.Context, id uuid.UUID, req UpdateProductParam) (Product, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

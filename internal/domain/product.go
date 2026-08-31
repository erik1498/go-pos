package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID             uuid.UUID       `json:"id"`
	CategoryID     uuid.UUID       `json:"category_id"`
	Category       Category        `json:"category,omitempty"`
	Name           string          `json:"name"`
	SKU            string          `json:"sku"`
	Price          decimal.Decimal `json:"price"`
	Stock          decimal.Decimal `json:"stock"`
	Taxes          []Tax           `json:"tax"`
	IdempotencyKey string          `json:"idempotency_key"`
	CreatedBy      uuid.UUID       `json:"created_by"`
	UpdatedBy      uuid.UUID       `json:"updated_by"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      *time.Time      `json:"deleted_at"`
}

type CreateProductParam struct {
	CategoryID uuid.UUID       `json:"category_id"`
	Name       string          `json:"name"`
	SKU        string          `json:"sku"`
	Price      decimal.Decimal `json:"price"`
	Taxes      []Tax           `json:"tax"`
}

type UpdateProductParam struct {
	CategoryID uuid.UUID       `json:"category_id"`
	Name       string          `json:"name"`
	SKU        string          `json:"sku"`
	Price      decimal.Decimal `json:"price"`
	Taxes      []Tax           `json:"tax"`
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

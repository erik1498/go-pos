package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID             uuid.UUID  `json:"id"`
	Name           string     `json:"name"`
	IdempotencyKey string     `json:"idempotency_key"`
	CreatedBy      uuid.UUID  `json:"created_by"`
	UpdatedBy      uuid.UUID  `json:"updated_by"`
	DeletedBy      *uuid.UUID `json:"deleted_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"-"`
}

type CreateCategoryParam struct {
	Name string `json:"name"`
}

type UpdateCategoryParam struct {
	Name string `json:"name"`
}

type CategoryRepository interface {
	Create(ctx context.Context, category Category) (Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (Category, error)
	GetAll(ctx context.Context, opts QueryOptions) ([]Category, int64, error)
	UpdateCategoryByID(ctx context.Context, id uuid.UUID, category Category) (Category, error)
	DeleteCategoryByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
}

type CategoryUsecase interface {
	Create(ctx context.Context, req CreateCategoryParam) (Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (Category, error)
	GetAll(ctx context.Context, opts QueryOptions) ([]Category, int64, error)
	UpdateCategoryByID(ctx context.Context, id uuid.UUID, req UpdateCategoryParam) (Category, error)
	DeleteCategoryByID(ctx context.Context, id uuid.UUID) error
}

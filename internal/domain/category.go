package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID             uuid.UUID
	Name           string
	IdempotencyKey string
	CreatedBy      uuid.UUID
	UpdatedBy      uuid.UUID
	DeletedBy      *uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

type CreateCategoryParam struct {
	Name string
}

type UpdateCategoryParam struct {
	Name string
}

type CategoryRepository interface {
	Create(ctx context.Context, category Category) (Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (Category, error)
	GetAll(ctx context.Context, opts QueryOptions) ([]Category, int64, error)
	UpdateByID(ctx context.Context, id uuid.UUID, category Category) (Category, error)
	DeleteByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
}

type CategoryUsecase interface {
	Create(ctx context.Context, req CreateCategoryParam) (Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (Category, error)
	GetAll(ctx context.Context, opts QueryOptions) ([]Category, int64, error)
	UpdateByID(ctx context.Context, id uuid.UUID, req UpdateCategoryParam) (Category, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

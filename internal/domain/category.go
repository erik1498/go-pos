package domain

import (
	"context"
	"errors"
	"go-pos/internal/model"

	"github.com/google/uuid"
)

var (
	ErrCategoryNotFound = errors.New("category data not found")
)

type CategoryRepository interface {
	Create(ctx context.Context, category model.Category) (model.Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Category, error)
	GetAll(ctx context.Context, opts QueryOptions) ([]model.Category, int64, error)
	UpdateCategoryByID(ctx context.Context, id uuid.UUID, category model.Category) (model.Category, error)
	DeleteCategoryByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
}

type CategoryUsecase interface {
	Create(ctx context.Context, req model.CategoryRequest) (model.Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Category, error)
	GetAll(ctx context.Context, opts QueryOptions) ([]model.Category, int64, error)
	UpdateCategoryByID(ctx context.Context, id uuid.UUID, req model.CategoryRequest) (model.Category, error)
	DeleteCategoryByID(ctx context.Context, id uuid.UUID) error
}

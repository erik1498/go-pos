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
	Create(category model.Category) (model.Category, error)
	GetByID(id uuid.UUID) (model.Category, error)
	GetAll(opts QueryOptions) ([]model.Category, int64, error)
	UpdateCategoryByID(id uuid.UUID, category model.Category) (model.Category, error)
	DeleteCategoryByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
}

type CategoryUsecase interface {
	Create(ctx context.Context, req model.Category) (model.Category, error)
	GetByID(id uuid.UUID) (model.Category, error)
	GetAll(opts QueryOptions) ([]model.Category, int64, error)
	UpdateCategoryByID(ctx context.Context, id uuid.UUID, req model.CategoryRequest) (model.Category, error)
	DeleteCategoryByID(ctx context.Context, id uuid.UUID) error
}

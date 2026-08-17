package domain

import (
	"errors"
	"go-pos/internal/model"

	"github.com/google/uuid"
)

var (
	CategoryErrNotFound = errors.New("category data not found")
)

type CategoryRepository interface {
	Create(category model.Category) (model.Category, error)
	GetByID(id uuid.UUID) (model.Category, error)
	GetAll(opts QueryOptions) ([]model.Category, int64, error)
	UpdateCategoryByID(id uuid.UUID, category model.Category) (model.Category, error)
	DeleteCategoryByID(id uuid.UUID) error
}

type CategoryUsecase interface {
	Create(category model.Category) (model.Category, error)
	GetByID(id uuid.UUID) (model.Category, error)
	GetAll(opts QueryOptions) ([]model.Category, int64, error)
	UpdateCategoryByID(id uuid.UUID, req model.Category) (model.Category, error)
	DeleteCategoryByID(id uuid.UUID) error
}

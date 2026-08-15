package category

import (
	"go-pos/internal/model"

	"github.com/google/uuid"
)

type Repository interface {
	GetCategoryList() ([]model.Category, error)
	CreateCategory(category model.Category) (model.Category, error)
	GetCategoryById(id uuid.UUID) (model.Category, error)
	UpdateCategoryById(category model.Category, id uuid.UUID) (model.Category, error)
	DeleteCategoryByID(id uuid.UUID) error
}

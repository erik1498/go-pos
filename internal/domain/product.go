package domain

import (
	"errors"
	"go-pos/internal/model"

	"github.com/google/uuid"
)

var (
	ErrProductNotFound               = errors.New("product data not found")
	ErrProductSKUIsAlreadyRegistered = errors.New("product sku is already registered")
)

type ProductRepository interface {
	GetAll(opts QueryOptions) ([]model.Product, int64, error)
	Create(product model.Product) (model.Product, error)
	GetByID(id uuid.UUID) (model.Product, error)
	UpdateByID(id uuid.UUID, product model.Product) (model.Product, error)
	DeleteByID(id uuid.UUID) error
}

type ProductUsecase interface {
	GetAll(opts QueryOptions) ([]model.Product, int64, error)
	Create(req model.ProductRequest) (model.Product, error)
	GetByID(id uuid.UUID) (model.Product, error)
	UpdateByID(id uuid.UUID, req model.ProductRequest) (model.Product, error)
	DeleteByID(id uuid.UUID) error
}

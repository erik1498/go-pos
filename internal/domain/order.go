package domain

import (
	"errors"
	"go-pos/internal/model"

	"github.com/google/uuid"
)

var (
	ErrOrderNotFound = errors.New("order not found")

	ErrOrderNoIsAlreadyRegistered = errors.New("order no is already registered")
)

type OrderRepository interface {
	GetAll(opts QueryOptions) ([]model.Order, int64, error)
	Create(order model.Order) (model.Order, error)
	GetByID(id uuid.UUID) (model.Order, error)
	UpdateByID(id uuid.UUID, order model.Order) (model.Order, error)
	DeleteByID(id uuid.UUID) error
}

type OrderUsecase interface {
	GetAll(opts QueryOptions) ([]model.Order, int64, error)
	Create(req model.OrderRequest) (model.Order, error)
}

package domain

import (
	"context"
	"errors"
	"go-pos/internal/model"

	"github.com/google/uuid"
)

var (
	ErrOrderNotFound = errors.New("order not found")

	ErrOrderNoIsAlreadyRegistered = errors.New("order no is already registered")
	ErrProductStockIsNotEnough    = errors.New("product stock is not enough")
)

type OrderRepository interface {
	GetAll(ctx context.Context, opts QueryOptions) ([]model.Order, int64, error)
	Create(ctx context.Context, order model.Order) (model.Order, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Order, error)
	UpdateByID(ctx context.Context, id uuid.UUID, order model.Order) (model.Order, error)
	DeleteByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
}

type OrderUsecase interface {
	GetAll(ctx context.Context, opts QueryOptions) ([]model.Order, int64, error)
	Create(ctx context.Context, req model.CreateOrderRequest) (model.Order, error)
}

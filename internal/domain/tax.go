package domain

import (
	"errors"
	"go-pos/internal/model"

	"github.com/google/uuid"
)

var (
	ErrTaxNotFound = errors.New("tax data not found")
)

type TaxRepository interface {
	GetAll(opts QueryOptions) ([]model.Tax, int64, error)
	Create(tax model.Tax) (model.Tax, error)
	GetByID(id uuid.UUID) (model.Tax, error)
	UpdateByID(id uuid.UUID, tax model.Tax) (model.Tax, error)
	DeleteByID(id uuid.UUID) error
}

type TaxUsecase interface {
	GetAll(opts QueryOptions) ([]model.Tax, int64, error)
	Create(req model.TaxRequest) (model.Tax, error)
	GetByID(id uuid.UUID) (model.Tax, error)
	UpdateByID(id uuid.UUID, tax model.TaxRequest) (model.Tax, error)
	DeleteByID(id uuid.UUID) error
}

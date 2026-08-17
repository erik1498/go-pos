package rest

import (
	"go-pos/internal/domain"
)

type handler struct {
	cUsecase domain.CategoryUsecase
	pUsecase domain.ProductUsecase
}

func NewHandler(
	cUsecase domain.CategoryUsecase,
	pUsecase domain.ProductUsecase,
) *handler {
	return &handler{
		cUsecase: cUsecase,
		pUsecase: pUsecase,
	}
}

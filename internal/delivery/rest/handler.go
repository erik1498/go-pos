package rest

import (
	"go-pos/internal/domain"
)

type handler struct {
	cUsecase domain.CategoryUsecase
}

func NewHandler(
	cUsecase domain.CategoryUsecase,
) *handler {
	return &handler{
		cUsecase: cUsecase,
	}
}

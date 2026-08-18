package rest

import (
	"go-pos/internal/domain"
)

type handler struct {
	cUsecase domain.CategoryUsecase
	pUsecase domain.ProductUsecase
	mUsecase domain.MemberUsecase
	oUsecase domain.OrderUsecase
}

func NewHandler(
	cUsecase domain.CategoryUsecase,
	pUsecase domain.ProductUsecase,
	mUsecase domain.MemberUsecase,
	oUsecase domain.OrderUsecase,
) *handler {
	return &handler{
		cUsecase: cUsecase,
		pUsecase: pUsecase,
		mUsecase: mUsecase,
		oUsecase: oUsecase,
	}
}

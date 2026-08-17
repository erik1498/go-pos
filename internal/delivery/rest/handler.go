package rest

import (
	"go-pos/internal/domain"
)

type handler struct {
	cUsecase domain.CategoryUsecase
	pUsecase domain.ProductUsecase
	mUsecase domain.MemberUsecase
}

func NewHandler(
	cUsecase domain.CategoryUsecase,
	pUsecase domain.ProductUsecase,
	mUsecase domain.MemberUsecase,
) *handler {
	return &handler{
		cUsecase: cUsecase,
		pUsecase: pUsecase,
		mUsecase: mUsecase,
	}
}

package rest

import (
	"encoding/json"
	"errors"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/response"
	"go-pos/pkg/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *handler) GetAllProduct(c echo.Context) error {
	opts := utils.ExtractQueryOptions(c)

	ctx := c.Request().Context()

	productList, totalItems, err := h.pUsecase.GetAll(ctx, opts)
	if err != nil {
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	meta := utils.BuildMetaPage(opts.Page, opts.Limit, totalItems)

	return response.SuccessWithMeta(c, http.StatusOK, domain.SuccessGetData, productList, meta)
}

func (h *handler) CreateProduct(c echo.Context) error {
	var req model.ProductRequest

	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrBadRequest.Error())
	}

	ctx := c.Request().Context()

	product, err := h.pUsecase.Create(ctx, req)
	if err != nil {
		if err == domain.ErrCategoryNotFound {
			return response.ErrNotFound(c, domain.ErrCategoryNotFound.Error())
		}
		if errors.Is(err, domain.ErrTaxNotFound) {
			return response.ErrNotFound(c, domain.ErrTaxNotFound.Error())
		}
		if errors.Is(err, domain.ErrProductSKUIsAlreadyRegistered) {
			return response.ErrBadRequest(c, domain.ErrProductSKUIsAlreadyRegistered.Error())
		}
		if errors.Is(err, domain.ErrIdempotencyKeyDuplicate) {
			return response.ErrConflictRequest(c, domain.ErrIdempotencyKeyDuplicate.Error())
		}
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusCreated, domain.SuccessCreateData, product)
}

func (h *handler) GetProductByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrIDInvalid.Error())
	}

	ctx := c.Request().Context()

	product, err := h.pUsecase.GetByID(ctx, ID)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return response.ErrNotFound(c, domain.ErrProductNotFound.Error())
		}
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusOK, domain.SuccessGetDataByID, product)
}

func (h *handler) UpdateProductByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrIDInvalid.Error())
	}

	var req model.ProductRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return response.ErrBadRequest(c, domain.ErrBadRequest.Error())
	}

	ctx := c.Request().Context()

	product, err := h.pUsecase.UpdateByID(ctx, ID, req)
	if err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			return response.ErrNotFound(c, domain.ErrCategoryNotFound.Error())
		}
		if errors.Is(err, domain.ErrProductNotFound) {
			return response.ErrNotFound(c, domain.ErrProductNotFound.Error())
		}
		if errors.Is(err, domain.ErrProductSKUIsAlreadyRegistered) {
			return response.ErrNotFound(c, domain.ErrProductSKUIsAlreadyRegistered.Error())
		}
		if errors.Is(err, domain.ErrTaxNotFound) {
			return response.ErrNotFound(c, domain.ErrTaxNotFound.Error())
		}
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusOK, domain.SuccessUpdateData, product)
}

func (h *handler) DeleteProductByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrBadRequest.Error())
	}

	ctx := c.Request().Context()

	err = h.pUsecase.DeleteByID(ctx, ID)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return response.ErrNotFound(c, domain.ErrProductNotFound.Error())
		}
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.NoContent(c)
}

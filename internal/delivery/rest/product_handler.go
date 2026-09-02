package rest

import (
	"fmt"
	"go-pos/internal/domain"
	"go-pos/pkg/response"
	"go-pos/pkg/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

type ProductTaxRequest struct {
	ID string `json:"id" validate:"required"`
}

type ProductRequest struct {
	Name       string           `json:"name" validate:"required"`
	SKU        string           `json:"sku" validate:"required"`
	Price      string           `json:"price" validate:"required"`
	CategoryID string           `json:"category_id" validate:"required"`
	Taxes      []ProductRequest `json:"taxes" validate:"required"`
}

type ProductResponse struct {
	ID             uuid.UUID        `json:"id"`
	CategoryID     uuid.UUID        `json:"category_id"`
	Category       CategoryResponse `json:"category"`
	Name           string           `json:"name"`
	SKU            string           `json:"sku"`
	Price          decimal.Decimal  `json:"price"`
	Stock          decimal.Decimal  `json:"stok"`
	Taxes          []TaxResponse    `json:"taxes"`
	IdempotencyKey string           `json:"idempotency_key"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

func toProductResponse(p domain.Product) ProductResponse {
	return ProductResponse{
		ID:             p.ID,
		CategoryID:     p.CategoryID,
		Name:           p.Name,
		SKU:            p.SKU,
		Price:          p.Price,
		Stock:          p.Stock,
		IdempotencyKey: p.IdempotencyKey,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}

func toProductResponseList(products []domain.Product) []ProductResponse {
	var res []ProductResponse
	for _, p := range products {
		res = append(res, toProductResponse(p))
	}
	return res
}

func (h *handler) GetAllProduct(c echo.Context) error {
	opts := utils.ExtractQueryOptions(c)

	ctx := c.Request().Context()

	products, totalItems, err := h.pUsecase.GetAll(ctx, opts)
	if err != nil {
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	meta := utils.BuildMetaPage(opts.Page, opts.Limit, totalItems)

	return response.SuccessWithMeta(c, http.StatusOK, domain.SuccessGetData, toProductResponseList(products), meta)
}

func (h *handler) CreateProduct(c echo.Context) error {
	var req ProductRequest
	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("[delivery][rest][product_handler][CreateProduct] invalid body: %w", domain.ErrBadRequest)
	}

	if err := c.Validate(&req); err != nil {
		return fmt.Errorf("[delivery][rest][product_handler][CreateProduct] validation error: %w", domain.ErrBadRequest)
	}

	ctx := c.Request().Context()

	param := domain.CreateProductParam{}

	product, err := h.pUsecase.Create(ctx, param)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusCreated, domain.SuccessCreateData, product)
}

func (h *handler) GetProductByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return fmt.Errorf("[delivery][rest][product_handler][GetProductByID] invalid UUID format: %w", domain.ErrIDInvalid)
	}

	ctx := c.Request().Context()

	product, err := h.pUsecase.GetByID(ctx, ID)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, domain.SuccessGetDataByID, product)
}

func (h *handler) UpdateProductByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return fmt.Errorf("[delivery][rest][product_handler][UpdateProductByID] invalid UUID format: %w", domain.ErrIDInvalid)
	}

	var req ProductRequest
	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("[delivery][rest][product_handler][UpdateProductByID] invalid body: %w", domain.ErrBadRequest)
	}

	if err := c.Validate(&req); err != nil {
		return fmt.Errorf("[delivery][rest][product_handler][UpdateProductByID] validation error: %w", domain.ErrBadRequest)
	}

	ctx := c.Request().Context()
	param := domain.UpdateProductParam{}

	product, err := h.pUsecase.UpdateByID(ctx, ID, param)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, domain.SuccessUpdateData, product)
}

func (h *handler) DeleteProductByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return fmt.Errorf("[delivery][rest][product_handler][DeleteProductByID] invalid UUID format: %w", domain.ErrIDInvalid)
	}

	ctx := c.Request().Context()

	err = h.pUsecase.DeleteByID(ctx, ID)
	if err != nil {
		return err
	}

	return response.NoContent(c)
}

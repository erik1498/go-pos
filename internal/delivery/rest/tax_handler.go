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

type TaxRequest struct {
	Name string          `json:"name" validate:"required,max=100"`
	Rate decimal.Decimal `json:"rate" validate:"required"`
}

type TaxResponse struct {
	Name      string
	Rate      decimal.Decimal
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toTaxResponse(t domain.Tax) TaxResponse {
	return TaxResponse{
		Name: t.Name,
		Rate: t.Rate,
	}
}

func toTaxResponseList(taxes []domain.Tax) []TaxResponse {
	var res []TaxResponse
	for _, t := range taxes {
		res = append(res, toTaxResponse(t))
	}
	return res
}

func (h *handler) GetAllTax(c echo.Context) error {
	opts := utils.ExtractQueryOptions(c)

	ctx := c.Request().Context()

	taxes, totalItems, err := h.tUsecase.GetAll(ctx, opts)
	if err != nil {
		return err
	}

	meta := utils.BuildMetaPage(opts.Page, opts.Limit, totalItems)

	return response.SuccessWithMeta(c, http.StatusOK, domain.SuccessGetData, toTaxResponseList(taxes), meta)
}

func (h *handler) CreateTax(c echo.Context) error {
	var req TaxRequest
	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("[delivery][rest][tax_handler][CreateTax] invalid body: %w", domain.ErrBadRequest)
	}

	if err := c.Validate(&req); err != nil {
		return fmt.Errorf("[delivery][rest][tax_handler][CreateTax] validation error: %w", domain.ErrBadRequest)
	}

	ctx := c.Request().Context()

	param := domain.CreateTaxParam{
		Name: req.Name,
		Rate: req.Rate,
	}
	tax, err := h.tUsecase.Create(ctx, param)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusCreated, domain.SuccessCreateData, toTaxResponse(tax))
}

func (h *handler) GetTaxByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return fmt.Errorf("[delivery][rest][tax_handler][GetTaxByID] invalid UUID format: %w", domain.ErrIDInvalid)
	}

	ctx := c.Request().Context()

	tax, err := h.tUsecase.GetByID(ctx, ID)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, domain.SuccessGetDataByID, toTaxResponse(tax))
}

func (h *handler) UpdateTaxByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return fmt.Errorf("[delivery][rest][tax_handler][UpdateTaxByID] invalid UUID format: %w", domain.ErrIDInvalid)
	}

	var req TaxRequest

	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("[delivery][rest][tax_handler][UpdateTaxByID] invalid body: %w", domain.ErrBadRequest)
	}

	if err := c.Validate(&req); err != nil {
		return fmt.Errorf("[delivery][rest][tax_handler][UpdateTaxByID] validation error: %w", domain.ErrBadRequest)
	}

	ctx := c.Request().Context()

	param := domain.UpdateTaxParam{
		Name: req.Name,
		Rate: req.Rate,
	}
	tax, err := h.tUsecase.UpdateByID(ctx, ID, param)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, domain.SuccessUpdateData, toTaxResponse(tax))
}

func (h *handler) DeleteTaxByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return fmt.Errorf("[delivery][rest][tax_handler][DeleteTaxByID] invalid UUID format: %w", domain.ErrIDInvalid)
	}

	ctx := c.Request().Context()

	if err := h.tUsecase.DeleteByID(ctx, ID); err != nil {
		return err
	}

	return response.NoContent(c)
}

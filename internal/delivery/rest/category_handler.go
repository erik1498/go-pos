package rest

import (
	"fmt"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/response"
	"go-pos/pkg/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *handler) GetAllCategory(c echo.Context) error {
	opts := utils.ExtractQueryOptions(c)

	ctx := c.Request().Context()

	categoryList, totalItems, err := h.cUsecase.GetAll(ctx, opts)
	if err != nil {
		return err
	}

	meta := utils.BuildMetaPage(opts.Page, opts.Limit, totalItems)

	return response.SuccessWithMeta(c, http.StatusOK, domain.SuccessGetData, categoryList, meta)
}

func (h *handler) CreateCategory(c echo.Context) error {
	var req model.CategoryRequest
	err := c.Bind(&req)

	if err != nil {
		return fmt.Errorf("[delivery][rest][category_handler][CreateCategory] invalid body: %w", domain.ErrBadRequest)
	}

	ctx := c.Request().Context()

	category, err := h.cUsecase.Create(ctx, req)

	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, domain.SuccessCreateData, category)
}

func (h *handler) GetCategoryByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return fmt.Errorf("[delivery][rest][category_handler][GetCategoryByID] invalid UUID format: %w", domain.ErrIDInvalid)
	}

	ctx := c.Request().Context()

	category, err := h.cUsecase.GetByID(ctx, ID)

	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, domain.SuccessGetDataByID, category)
}

func (h *handler) UpdateCategoryById(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return fmt.Errorf("[delivery][rest][category_handler][UpdateCategoryById] invalid UUID format: %w", domain.ErrIDInvalid)
	}

	var req model.CategoryRequest
	err = c.Bind(&req)

	if err != nil {
		return fmt.Errorf("[delivery][rest][category_handler][UpdateCategoryById] invalid body: %w", domain.ErrBadRequest)
	}

	ctx := c.Request().Context()
	category, err := h.cUsecase.UpdateCategoryByID(ctx, ID, req)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, domain.SuccessUpdateData, category)
}

func (h *handler) DeleteCategoryByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return fmt.Errorf("[delivery][rest][category_handler][DeleteCategoryByID] invalid UUID format: %w", domain.ErrIDInvalid)
	}

	ctx := c.Request().Context()

	err = h.cUsecase.DeleteCategoryByID(ctx, ID)
	if err != nil {
		return err
	}

	return response.NoContent(c)
}

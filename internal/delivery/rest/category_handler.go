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
)

type CategoryRequest struct {
	Name string `json:"name" validate:"required,max=100"`
}

type CategoryResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toCategoryResponse(c domain.Category) CategoryResponse {
	return CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func toCategoryResponseList(categories []domain.Category) []CategoryResponse {
	var res []CategoryResponse
	for _, c := range categories {
		res = append(res, toCategoryResponse(c))
	}
	return res
}

func (h *handler) GetAllCategory(c echo.Context) error {
	opts := utils.ExtractQueryOptions(c)
	ctx := c.Request().Context()

	categories, totalItems, err := h.cUsecase.GetAll(ctx, opts)
	if err != nil {
		return err
	}

	meta := utils.BuildMetaPage(opts.Page, opts.Limit, totalItems)

	return response.SuccessWithMeta(c, http.StatusOK, domain.SuccessGetData, toCategoryResponseList(categories), meta)
}

func (h *handler) CreateCategory(c echo.Context) error {
	var req CategoryRequest
	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("[delivery][rest][category_handler][CreateCategory] invalid body: %w", domain.ErrBadRequest)
	}

	if err := c.Validate(&req); err != nil {
		return fmt.Errorf("[delivery][rest][category_handler][CreateCategory] validation error: %w", domain.ErrBadRequest)
	}

	ctx := c.Request().Context()

	param := domain.CreateCategoryParam{
		Name: req.Name,
	}
	category, err := h.cUsecase.Create(ctx, param)

	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, domain.SuccessCreateData, toCategoryResponse(category))
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

	return response.Success(c, http.StatusOK, domain.SuccessGetDataByID, toCategoryResponse(category))
}

func (h *handler) UpdateCategoryById(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return fmt.Errorf("[delivery][rest][category_handler][UpdateCategoryById] invalid UUID format: %w", domain.ErrIDInvalid)
	}

	var req CategoryRequest
	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("[delivery][rest][category_handler][UpdateCategoryById] invalid body: %w", domain.ErrBadRequest)
	}

	if err := c.Validate(&req); err != nil {
		return fmt.Errorf("[delivery][rest][category_handler][UpdateCategoryById] validation error: %w", domain.ErrBadRequest)
	}

	ctx := c.Request().Context()
	param := domain.UpdateCategoryParam{
		Name: req.Name,
	}
	category, err := h.cUsecase.UpdateByID(ctx, ID, param)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, domain.SuccessUpdateData, toCategoryResponse(category))
}

func (h *handler) DeleteCategoryByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return fmt.Errorf("[delivery][rest][category_handler][DeleteCategoryByID] invalid UUID format: %w", domain.ErrIDInvalid)
	}

	ctx := c.Request().Context()

	err = h.cUsecase.DeleteByID(ctx, ID)
	if err != nil {
		return err
	}

	return response.NoContent(c)
}

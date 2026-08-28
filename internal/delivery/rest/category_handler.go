package rest

import (
	"encoding/json"
	"errors"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/response"
	"go-pos/pkg/utils"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *handler) GetAllCategory(c echo.Context) error {
	opts := utils.ExtractQueryOptions(c)

	ctx := c.Request().Context()

	categoryList, totalItems, err := h.cUsecase.GetAll(ctx, opts)
	if err != nil {
		slog.Error("[delivery][rest][category_handler][GetAllCategory] failed to fetch data",
			slog.String("error", err.Error()),
		)
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	meta := utils.BuildMetaPage(opts.Page, opts.Limit, totalItems)

	return response.SuccessWithMeta(c, http.StatusOK, domain.SuccessGetData, categoryList, meta)
}

func (h *handler) CreateCategory(c echo.Context) error {
	var req model.CategoryRequest
	err := json.NewDecoder(c.Request().Body).Decode(&req)

	if err != nil {
		slog.Warn("[delivery][rest][category_handler][CreateCategory] invalid request body",
			slog.String("error", err.Error()),
		)
		return response.ErrBadRequest(c, domain.ErrBadRequest.Error())
	}

	ctx := c.Request().Context()

	category, err := h.cUsecase.Create(ctx, req)

	if err != nil {
		if errors.Is(err, domain.ErrIdempotencyKeyDuplicate) {
			slog.Warn("[delivery][rest][category_handler][CreateCategory] duplicate idempotency key",
				slog.String("error", err.Error()),
			)
			return response.ErrConflictRequest(c, domain.ErrIdempotencyKeyDuplicate.Error())
		}
		slog.Error("[delivery][rest][category_handler][CreateCategory] failed to create data",
			slog.String("error", err.Error()),
		)
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusOK, domain.SuccessCreateData, category)
}

func (h *handler) GetCategoryByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		slog.Warn("[delivery][rest][category_handler][GetCategoryByID] invalid id",
			slog.String("error", err.Error()),
		)
		return response.ErrBadRequest(c, domain.ErrIDInvalid.Error())
	}

	ctx := c.Request().Context()

	category, err := h.cUsecase.GetByID(ctx, ID)

	if err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			slog.Warn("[delivery][rest][category_handler][GetCategoryByID] duplicate idempotency key",
				slog.String("error", err.Error()),
			)
			return response.ErrNotFound(c, domain.ErrCategoryNotFound.Error())
		}
		slog.Error("[delivery][rest][category_handler][GetCategoryByID] failed to fetch data",
			slog.String("error", err.Error()),
		)
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusOK, domain.SuccessGetDataByID, category)
}

func (h *handler) UpdateCategoryById(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		slog.Warn("[delivery][rest][category_handler][UpdateCategoryById] invalid id",
			slog.String("error", err.Error()),
		)
		return response.ErrBadRequest(c, domain.ErrIDInvalid.Error())
	}

	var req model.CategoryRequest

	err = json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		slog.Warn("[delivery][rest][category_handler][UpdateCategoryById] invalid request body",
			slog.String("error", err.Error()),
		)
		return response.ErrBadRequest(c, domain.ErrBadRequest.Error())
	}

	ctx := c.Request().Context()
	category, err := h.cUsecase.UpdateCategoryByID(ctx, ID, req)
	if err != nil {

		if errors.Is(err, domain.ErrCategoryNotFound) {
			slog.Warn("[delivery][rest][category_handler][UpdateCategoryById] duplicate idempotency key",
				slog.String("error", err.Error()),
			)
			return response.ErrNotFound(c, domain.ErrCategoryNotFound.Error())
		}
		slog.Error("[delivery][rest][category_handler][UpdateCategoryById] failed to update data",
			slog.String("error", err.Error()),
		)
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusOK, domain.SuccessUpdateData, category)
}

func (h *handler) DeleteCategoryByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		slog.Warn("[delivery][rest][category_handler][DeleteCategoryByID] invalid id",
			slog.String("error", err.Error()),
		)
		return response.ErrBadRequest(c, domain.ErrIDInvalid.Error())
	}

	ctx := c.Request().Context()

	err = h.cUsecase.DeleteCategoryByID(ctx, ID)
	if err != nil {

		if errors.Is(err, domain.ErrCategoryNotFound) {
			slog.Warn("[delivery][rest][category_handler][DeleteCategoryByID] duplicate idempotency key",
				slog.String("error", err.Error()),
			)
			return response.ErrNotFound(c, domain.ErrCategoryNotFound.Error())
		}

		slog.Error("[delivery][rest][category_handler][DeleteCategoryByID] failed to delete data",
			slog.String("error", err.Error()),
		)
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.NoContent(c)
}

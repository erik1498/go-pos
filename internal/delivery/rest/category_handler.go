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
	"gorm.io/gorm"
)

func (h *handler) GetAll(c echo.Context) error {
	opts := utils.ExtractQueryOptions(c)

	categoryList, totalItems, err := h.cUsecase.GetAll(opts)
	if err != nil {
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	meta := utils.BuildMetaPage(opts.Page, opts.Limit, totalItems)

	return response.SuccessWithMeta(c, http.StatusOK, domain.SuccessGetData, categoryList, meta)
}

func (h *handler) Create(c echo.Context) error {
	var req model.Category
	err := json.NewDecoder(c.Request().Body).Decode(&req)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	category, err := h.cUsecase.Create(req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": category,
	})
}

func (h *handler) GetByPublicID(c echo.Context) error {
	idParam := c.Param("id")

	publicID, err := uuid.Parse(idParam)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrIDInvalid.Error())
	}

	category, err := h.cUsecase.GetByPublicID(publicID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, err.Error())
		}

		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusOK, domain.SuccessGetDataByID, category)
}

func (h *handler) UpdateCategoryById(c echo.Context) error {
	idParam := c.Param("id")

	publicId, err := uuid.Parse(idParam)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrIDInvalid.Error())
	}

	var req model.Category

	err = json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrBadRequest.Error())
	}

	category, err := h.cUsecase.UpdateCategoryByID(publicId, req)
	if err != nil {

		if errors.Is(err, domain.CategoryErrNotFound) {
			return response.ErrNotFound(c, domain.CategoryErrNotFound.Error())
		}

		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusOK, domain.SuccessUpdateData, category)
}

func (h *handler) DeleteCategoryByID(c echo.Context) error {
	idParam := c.Param("id")

	publicID, err := uuid.Parse(idParam)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrIDInvalid.Error())
	}

	err = h.cUsecase.DeleteCategoryByID(publicID)
	if err != nil {

		if errors.Is(err, domain.CategoryErrNotFound) {
			return response.ErrNotFound(c, domain.CategoryErrNotFound.Error())
		}

		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.NoContent(c)
}

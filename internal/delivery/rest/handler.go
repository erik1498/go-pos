package rest

import (
	"encoding/json"
	"go-pos/internal/model"
	"go-pos/internal/usecase/pos"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type handler struct {
	pUsecase pos.Usecase
}

func NewHandler(
	pUsecase pos.Usecase,
) *handler {
	return &handler{
		pUsecase: pUsecase,
	}
}

func (h *handler) GetCategoryList(c echo.Context) error {
	categoryList, err := h.pUsecase.GetCategoryList()

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": categoryList,
	})
}

func (h *handler) CreateCategory(c echo.Context) error {
	var req model.Category
	err := json.NewDecoder(c.Request().Body).Decode(&req)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	category, err := h.pUsecase.CreateCategory(req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": category,
	})
}

func (h *handler) GetCategoryById(c echo.Context) error {
	idParam := c.Param("id")

	publicID, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "format ID tidak valid",
		})
	}

	category, err := h.pUsecase.GetCategoryById(publicID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": category,
	})
}

func (h *handler) UpdateCategoryById(c echo.Context) error {
	idParam := c.Param("id")

	publicId, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "format ID tidak valid",
		})
	}

	var req model.Category

	err = json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	category, err := h.pUsecase.UpdateCategoryById(req, publicId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": category,
	})
}

func (h *handler) DeleteCategoryByID(c echo.Context) error {
	idParam := c.Param("id")

	publicID, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "format ID tidak valid",
		})
	}

	err = h.pUsecase.DeleteCategoryByID(publicID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, map[string]interface{}{})
}

package rest

import "github.com/labstack/echo/v4"

func LoadRoutes(e *echo.Echo, h *handler) {
	categoryGroup := e.Group("/categories")

	categoryGroup.GET("", h.GetCategoryList)
	categoryGroup.POST("", h.CreateCategory)
	categoryGroup.GET("/:id", h.GetCategoryById)
	categoryGroup.PUT("/:id", h.UpdateCategoryById)
	categoryGroup.DELETE("/:id", h.DeleteCategoryByID)
}

package rest

import "github.com/labstack/echo/v4"

func LoadRoutes(e *echo.Echo, h *handler) {
	categoryGroup := e.Group("/categories")

	categoryGroup.GET("", h.GetAll)
	categoryGroup.POST("", h.Create)
	categoryGroup.GET("/:id", h.GetByPublicID)
	categoryGroup.PUT("/:id", h.UpdateCategoryById)
	categoryGroup.DELETE("/:id", h.DeleteCategoryByID)
}

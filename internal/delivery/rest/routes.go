package rest

import "github.com/labstack/echo/v4"

func LoadRoutes(e *echo.Echo, h *handler) {
	categoryGroup := e.Group("/categories")
	categoryGroup.GET("", h.GetAllCategory)
	categoryGroup.POST("", h.CreateCategory)
	categoryGroup.GET("/:id", h.GetCategoryByID)
	categoryGroup.PUT("/:id", h.UpdateCategoryById)
	categoryGroup.DELETE("/:id", h.DeleteCategoryByID)

	productGroup := e.Group("/products")
	productGroup.GET("", h.GetAllProduct)
	productGroup.POST("", h.CreateProduct)
	productGroup.GET("/:id", h.GetProductByID)
	productGroup.PUT("/:id", h.UpdateProductByID)
	productGroup.DELETE("/:id", h.DeleteProductByID)

	memberGroup := e.Group("/members")
	memberGroup.GET("", h.GetAllMember)
	memberGroup.POST("", h.CreateMember)
	memberGroup.GET("/:id", h.GetMemberByID)
	memberGroup.PUT("/:id", h.UpdateMemberByID)
	memberGroup.DELETE("/:id", h.DeleteMemberByID)

	orderGroup := e.Group("/orders")
	orderGroup.GET("", h.GetAllOrder)
	orderGroup.POST("", h.CreateOrder)

	taxGroup := e.Group("/tax")
	taxGroup.GET("", h.GetAllTax)
	taxGroup.POST("", h.CreateTax)
	taxGroup.GET("/:id", h.GetTaxByID)
	taxGroup.PUT("/:id", h.GetTaxByID)
	taxGroup.DELETE("/:id", h.DeleteTaxByID)
}

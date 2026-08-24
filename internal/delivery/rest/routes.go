package rest

import "github.com/labstack/echo/v4"

func LoadRoutes(e *echo.Echo, h *handler) {
	am := GetAuthMiddleware(h.uUsecase)

	v1 := e.Group("/api/v1")
	apiPublic := v1.Group("")

	userGroup := apiPublic.Group("/users")
	userGroup.POST("/register", h.RegisterUser)
	userGroup.POST("/login", h.LoginUser)

	apiPrivateGroup := v1.Group("")
	apiPrivateGroup.Use(am.CheckAuth)

	categoryGroup := apiPrivateGroup.Group("/categories")
	categoryGroup.GET("", h.GetAllCategory)
	categoryGroup.POST("", h.CreateCategory)
	categoryGroup.GET("/:id", h.GetCategoryByID)
	categoryGroup.PUT("/:id", h.UpdateCategoryById)
	categoryGroup.DELETE("/:id", h.DeleteCategoryByID)

	productGroup := apiPrivateGroup.Group("/products")
	productGroup.GET("", h.GetAllProduct)
	productGroup.POST("", h.CreateProduct)
	productGroup.GET("/:id", h.GetProductByID)
	productGroup.PUT("/:id", h.UpdateProductByID)
	productGroup.DELETE("/:id", h.DeleteProductByID)

	memberGroup := apiPrivateGroup.Group("/members")
	memberGroup.GET("", h.GetAllMember)
	memberGroup.POST("", h.CreateMember)
	memberGroup.GET("/:id", h.GetMemberByID)
	memberGroup.PUT("/:id", h.UpdateMemberByID)
	memberGroup.DELETE("/:id", h.DeleteMemberByID)

	orderGroup := apiPrivateGroup.Group("/orders")
	orderGroup.GET("", h.GetAllOrder)
	orderGroup.POST("", h.CreateOrder)

	taxGroup := apiPrivateGroup.Group("/taxes")
	taxGroup.GET("", h.GetAllTax)
	taxGroup.POST("", h.CreateTax)
	taxGroup.GET("/:id", h.GetTaxByID)
	taxGroup.PUT("/:id", h.UpdateTaxByID)
	taxGroup.DELETE("/:id", h.DeleteTaxByID)
}

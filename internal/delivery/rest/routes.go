package rest

import (
	"go-pos/internal/model"

	"github.com/labstack/echo/v4"
)

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
	categoryGroup.GET("", h.GetAllCategory, am.RequiredRoles(string(model.UserRoleAdmin)))
	categoryGroup.POST("", h.CreateCategory, am.RequiredRoles(string(model.UserRoleAdmin)))
	categoryGroup.GET("/:id", h.GetCategoryByID, am.RequiredRoles(string(model.UserRoleAdmin)))
	categoryGroup.PUT("/:id", h.UpdateCategoryById, am.RequiredRoles(string(model.UserRoleAdmin)))
	categoryGroup.DELETE("/:id", h.DeleteCategoryByID, am.RequiredRoles(string(model.UserRoleAdmin)))

	productGroup := apiPrivateGroup.Group("/products")
	productGroup.GET("", h.GetAllProduct, am.RequiredRoles(string(model.UserRoleAdmin)))
	productGroup.POST("", h.CreateProduct, am.RequiredRoles(string(model.UserRoleAdmin)))
	productGroup.GET("/:id", h.GetProductByID, am.RequiredRoles(string(model.UserRoleAdmin)))
	productGroup.PUT("/:id", h.UpdateProductByID, am.RequiredRoles(string(model.UserRoleAdmin)))
	productGroup.DELETE("/:id", h.DeleteProductByID, am.RequiredRoles(string(model.UserRoleAdmin)))

	memberGroup := apiPrivateGroup.Group("/members")
	memberGroup.GET("", h.GetAllMember, am.RequiredRoles(string(model.UserRoleAdmin)))
	memberGroup.POST("", h.CreateMember, am.RequiredRoles(string(model.UserRoleAdmin)))
	memberGroup.GET("/:id", h.GetMemberByID, am.RequiredRoles(string(model.UserRoleAdmin)))
	memberGroup.PUT("/:id", h.UpdateMemberByID, am.RequiredRoles(string(model.UserRoleAdmin)))
	memberGroup.DELETE("/:id", h.DeleteMemberByID, am.RequiredRoles(string(model.UserRoleAdmin)))

	orderGroup := apiPrivateGroup.Group("/orders")
	orderGroup.GET("", h.GetAllOrder, am.RequiredRoles(string(model.UserRoleAdmin)))
	orderGroup.POST("", h.CreateOrder, am.RequiredRoles(string(model.UserRoleAdmin)))

	taxGroup := apiPrivateGroup.Group("/taxes")
	taxGroup.GET("", h.GetAllTax, am.RequiredRoles(string(model.UserRoleAdmin)))
	taxGroup.POST("", h.CreateTax, am.RequiredRoles(string(model.UserRoleAdmin)))
	taxGroup.GET("/:id", h.GetTaxByID, am.RequiredRoles(string(model.UserRoleAdmin)))
	taxGroup.PUT("/:id", h.UpdateTaxByID, am.RequiredRoles(string(model.UserRoleAdmin)))
	taxGroup.DELETE("/:id", h.DeleteTaxByID, am.RequiredRoles(string(model.UserRoleAdmin)))
}

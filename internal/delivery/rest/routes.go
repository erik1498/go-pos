package rest

import (
	"go-pos/internal/model"
	"go-pos/pkg/middleware"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func LoadRoutes(e *echo.Echo, h *handler, rdb *redis.Client) {
	am := middleware.GetAuthMiddleware(h.uUsecase)

	v1 := e.Group("/api/v1")
	apiPublic := v1.Group("")

	userGroup := apiPublic.Group("/users")
	userGroup.POST("/register", h.RegisterUser)
	userGroup.POST("/login", h.LoginUser)

	apiPrivateGroup := v1.Group("")
	apiPrivateGroup.Use(am.CheckAuth)

	categoryGroup := apiPrivateGroup.Group("/categories")
	categoryGroup.GET("", h.GetAllCategory, am.RequiredRoles(string(model.UserRoleCashier)))
	categoryGroup.POST("", h.CreateCategory, am.RequiredRoles(string(model.UserRoleCashier)), middleware.Idempotency(rdb))
	categoryGroup.GET("/:id", h.GetCategoryByID, am.RequiredRoles(string(model.UserRoleCashier)))
	categoryGroup.PUT("/:id", h.UpdateCategoryById, am.RequiredRoles(string(model.UserRoleCashier)))
	categoryGroup.DELETE("/:id", h.DeleteCategoryByID, am.RequiredRoles(string(model.UserRoleCashier)))

	productGroup := apiPrivateGroup.Group("/products")
	productGroup.GET("", h.GetAllProduct, am.RequiredRoles(string(model.UserRoleCashier)))
	productGroup.POST("", h.CreateProduct, am.RequiredRoles(string(model.UserRoleCashier)))
	productGroup.GET("/:id", h.GetProductByID, am.RequiredRoles(string(model.UserRoleCashier)))
	productGroup.PUT("/:id", h.UpdateProductByID, am.RequiredRoles(string(model.UserRoleCashier)))
	productGroup.DELETE("/:id", h.DeleteProductByID, am.RequiredRoles(string(model.UserRoleCashier)))

	memberGroup := apiPrivateGroup.Group("/members")
	memberGroup.GET("", h.GetAllMember, am.RequiredRoles(string(model.UserRoleCashier)))
	memberGroup.POST("", h.CreateMember, am.RequiredRoles(string(model.UserRoleCashier)))
	memberGroup.GET("/:id", h.GetMemberByID, am.RequiredRoles(string(model.UserRoleCashier)))
	memberGroup.PUT("/:id", h.UpdateMemberByID, am.RequiredRoles(string(model.UserRoleCashier)))
	memberGroup.DELETE("/:id", h.DeleteMemberByID, am.RequiredRoles(string(model.UserRoleCashier)))

	orderGroup := apiPrivateGroup.Group("/orders")
	orderGroup.GET("", h.GetAllOrder, am.RequiredRoles(string(model.UserRoleCashier)))
	orderGroup.POST("", h.CreateOrder, am.RequiredRoles(string(model.UserRoleCashier)))

	taxGroup := apiPrivateGroup.Group("/taxes")
	taxGroup.GET("", h.GetAllTax, am.RequiredRoles(string(model.UserRoleCashier)))
	taxGroup.POST("", h.CreateTax, am.RequiredRoles(string(model.UserRoleAdmin)))
	taxGroup.GET("/:id", h.GetTaxByID, am.RequiredRoles(string(model.UserRoleAdmin)))
	taxGroup.PUT("/:id", h.UpdateTaxByID, am.RequiredRoles(string(model.UserRoleAdmin)))
	taxGroup.DELETE("/:id", h.DeleteTaxByID, am.RequiredRoles(string(model.UserRoleAdmin)))
}

package rest

import (
	"go-pos/internal/domain"
	"go-pos/pkg/middleware"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func LoadRoutes(e *echo.Echo, h *handler, rdb *redis.Client) {
	e.Use(middleware.RateLimiter(rdb, 100, 1*time.Minute))

	am := middleware.GetAuthMiddleware(h.uUsecase)

	v1 := e.Group("/api/v1")
	apiPublic := v1.Group("")

	userGroup := apiPublic.Group("/users")
	userGroup.POST("/register", h.RegisterUser, middleware.Idempotency(rdb, "user"))
	userGroup.POST("/login", h.LoginUser, middleware.RateLimiter(rdb, 5, 1*time.Minute))

	apiPrivateGroup := v1.Group("")
	apiPrivateGroup.Use(am.CheckAuth)

	categoryGroup := apiPrivateGroup.Group("/categories")
	categoryGroup.GET("", h.GetAllCategory, am.RequiredRoles(string(domain.UserRoleAdmin)))
	categoryGroup.POST("", h.CreateCategory, am.RequiredRoles(string(domain.UserRoleAdmin)), middleware.Idempotency(rdb, "ctgr"))
	categoryGroup.GET("/:id", h.GetCategoryByID, am.RequiredRoles(string(domain.UserRoleAdmin)))
	categoryGroup.PUT("/:id", h.UpdateCategoryById, am.RequiredRoles(string(domain.UserRoleAdmin)))
	categoryGroup.DELETE("/:id", h.DeleteCategoryByID, am.RequiredRoles(string(domain.UserRoleAdmin)))

	productGroup := apiPrivateGroup.Group("/products")
	productGroup.GET("", h.GetAllProduct, am.RequiredRoles(string(domain.UserRoleAdmin)))
	productGroup.POST("", h.CreateProduct, am.RequiredRoles(string(domain.UserRoleAdmin)), middleware.Idempotency(rdb, "prdc"))
	productGroup.GET("/:id", h.GetProductByID, am.RequiredRoles(string(domain.UserRoleAdmin)))
	productGroup.PUT("/:id", h.UpdateProductByID, am.RequiredRoles(string(domain.UserRoleAdmin)))
	productGroup.DELETE("/:id", h.DeleteProductByID, am.RequiredRoles(string(domain.UserRoleAdmin)))

	memberGroup := apiPrivateGroup.Group("/members")
	memberGroup.GET("", h.GetAllMember, am.RequiredRoles(string(domain.UserRoleAdmin)))
	memberGroup.POST("", h.CreateMember, am.RequiredRoles(string(domain.UserRoleAdmin)), middleware.Idempotency(rdb, "mmbr"))
	memberGroup.GET("/:id", h.GetMemberByID, am.RequiredRoles(string(domain.UserRoleAdmin)))
	memberGroup.PUT("/:id", h.UpdateMemberByID, am.RequiredRoles(string(domain.UserRoleAdmin)))

	orderGroup := apiPrivateGroup.Group("/orders")
	orderGroup.GET("", h.GetAllOrder, am.RequiredRoles(string(domain.UserRoleAdmin)))
	orderGroup.POST("", h.CreateOrder, am.RequiredRoles(string(domain.UserRoleAdmin)))

	taxGroup := apiPrivateGroup.Group("/taxes")
	taxGroup.GET("", h.GetAllTax, am.RequiredRoles(string(domain.UserRoleAdmin)))
	taxGroup.POST("", h.CreateTax, am.RequiredRoles(string(domain.UserRoleAdmin)), middleware.Idempotency(rdb, "taxs"))
	taxGroup.GET("/:id", h.GetTaxByID, am.RequiredRoles(string(domain.UserRoleAdmin)))
	taxGroup.PUT("/:id", h.UpdateTaxByID, am.RequiredRoles(string(domain.UserRoleAdmin)))
	taxGroup.DELETE("/:id", h.DeleteTaxByID, am.RequiredRoles(string(domain.UserRoleAdmin)))
}

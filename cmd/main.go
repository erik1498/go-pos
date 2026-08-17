package main

import (
	"go-pos/internal/database"
	"go-pos/internal/delivery/rest"
	"go-pos/internal/repository"
	"go-pos/internal/usecase"

	"github.com/labstack/echo/v4"
)

const (
	dbAddress = "host=localhost user=postgres password=postgres dbname=go-pos sslmode=disable"
)

func main() {
	e := echo.New()

	db := database.GetDB(dbAddress)

	cRepo := repository.GetCategoryRepository(db)
	pRepo := repository.GetProductRepository(db)

	cUsecase := usecase.GetCategoryUsecase(cRepo)
	pUsecase := usecase.GetProductUsecase(pRepo, cRepo)

	handler := rest.NewHandler(cUsecase, pUsecase)

	rest.LoadRoutes(e, handler)

	e.Logger.Fatal(e.Start(":3000"))
}

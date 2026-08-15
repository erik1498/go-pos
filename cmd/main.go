package main

import (
	"go-pos/internal/database"
	"go-pos/internal/delivery/rest"
	"go-pos/internal/repository/category"
	"go-pos/internal/usecase/pos"

	"github.com/labstack/echo/v4"
)

const (
	dbAddress = "host=localhost user=postgres password=postgres dbname=go-pos sslmode=disable"
)

func main() {
	e := echo.New()

	db := database.GetDB(dbAddress)

	cRepo := category.GetRepository(db)

	pUsecase := pos.GetUsecase(cRepo)

	handler := rest.NewHandler(pUsecase)

	rest.LoadRoutes(e, handler)

	e.Logger.Fatal(e.Start(":3000"))
}

package main

import (
	"go-pos/internal/database"
	"go-pos/internal/delivery/rest"
	"go-pos/internal/domain"
	"go-pos/internal/repository"
	"go-pos/internal/usecase"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

const (
	dbAddress = "host=localhost user=postgres password=postgres dbname=go-pos sslmode=disable"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(domain.AlertEnvNotFound)
	}

	aesKey := os.Getenv("AES_256_KEY")
	bindexKey := os.Getenv("BLIND_INDEX_KEY")

	if len(aesKey) != 32 {
		panic(domain.AlertAESNot32Character)
	}

	if bindexKey == "" {
		panic(domain.AlertBlindIndexEmpty)
	}

	e := echo.New()

	db := database.GetDB(dbAddress)

	cRepo := repository.GetCategoryRepository(db)
	pRepo := repository.GetProductRepository(db)
	mRepo := repository.GetMemberRepository(db)
	oRepo := repository.GetOrderRepository(db)

	cUsecase := usecase.GetCategoryUsecase(cRepo)
	pUsecase := usecase.GetProductUsecase(pRepo, cRepo)
	mUsecase := usecase.GetMemberUsecase(mRepo, aesKey, bindexKey)
	oUsecase := usecase.GetOrderUsecase(oRepo, mRepo)

	handler := rest.NewHandler(cUsecase, pUsecase, mUsecase, oUsecase)

	rest.LoadRoutes(e, handler)

	e.Logger.Fatal(e.Start(":3000"))
}

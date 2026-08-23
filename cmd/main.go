package main

import (
	"go-pos/internal/database"
	"go-pos/internal/delivery/rest"
	"go-pos/internal/domain"
	"go-pos/internal/repository"
	"go-pos/internal/usecase"
	"go-pos/pkg/utils"
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

	rsaPrivateKey, err := utils.LoadRSAPrivateKey()
	if err != nil {
		panic(err)
	}

	rsaPublicKey, err := utils.LoadRSAPublicKey()
	if err != nil {
		panic(err)
	}

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
	tRepo := repository.GetTaxRepository(db)
	uRepo := repository.GetUserRepository(db)

	cUsecase := usecase.GetCategoryUsecase(cRepo)
	pUsecase := usecase.GetProductUsecase(pRepo, cRepo, tRepo)
	mUsecase := usecase.GetMemberUsecase(mRepo, aesKey, bindexKey)
	oUsecase := usecase.GetOrderUsecase(oRepo, mRepo, pRepo)
	tUsecase := usecase.GetTaxUsecase(tRepo)
	uUsecase := usecase.GetUserUsecase(uRepo, aesKey, bindexKey, rsaPrivateKey, rsaPublicKey)

	handler := rest.NewHandler(cUsecase, pUsecase, mUsecase, oUsecase, tUsecase, uUsecase)

	rest.LoadRoutes(e, handler)

	e.Logger.Fatal(e.Start(":3000"))
}

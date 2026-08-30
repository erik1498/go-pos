package main

import (
	"context"
	"go-pos/internal/database"
	"go-pos/internal/delivery/rest"
	"go-pos/internal/domain"
	"go-pos/internal/repository"
	"go-pos/internal/usecase"
	"go-pos/pkg/middleware"
	"go-pos/pkg/response"
	"go-pos/pkg/utils"
	"log/slog" // Beralih penuh ke slog
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMid "github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

func initRedis() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:         host + ":" + port,
		Password:     password,
		DB:           0,
		PoolSize:     20,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		slog.Error("REDIS: CONNECTION FAILED", slog.String("host", host), slog.String("port", port), slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("REDIS: CONNECTED")
	return client
}

func main() {
	middleware.InitSlog()

	if err := godotenv.Load(); err != nil {
		slog.Info(".env not found. Using system environment variables.")
	}

	aesKey := os.Getenv("AES_256_KEY")
	bindexKey := os.Getenv("BLIND_INDEX_KEY")

	if len(aesKey) != 32 {
		slog.Error("SECURITY", slog.String("error", domain.AlertAESNot32Character))
		os.Exit(1)
	}
	if bindexKey == "" {
		slog.Error("SECURITY", slog.String("error", domain.AlertBlindIndexEmpty))
		os.Exit(1)
	}

	rsaPrivateKey, err := utils.LoadRSAPrivateKey()
	if err != nil {
		slog.Error("RSA: LOAD PRIVATE KEY FAILED", slog.String("error", err.Error()))
		os.Exit(1)
	}

	rsaPublicKey, err := utils.LoadRSAPublicKey()
	if err != nil {
		slog.Error("RSA: LOAD PUBLIC KEY FAILED", slog.String("error", err.Error()))
		os.Exit(1)
	}

	dbAddress := os.Getenv("DATABASE_URL")
	if dbAddress == "" {
		dbAddress = "host=localhost user=postgres password=postgres dbname=go-pos sslmode=disable"
	}
	db := database.GetDB(dbAddress)
	rdb := initRedis()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.HTTPErrorHandler = response.CustomHTTPErrorHandler

	e.Use(echoMid.Recover())
	e.Use(echoMid.RequestID())
	e.Use(middleware.SlogMiddleware())

	cRepo := repository.GetCategoryRepository(db)
	pRepo := repository.GetProductRepository(db)
	mRepo := repository.GetMemberRepository(db)
	oRepo := repository.GetOrderRepository(db)
	tRepo := repository.GetTaxRepository(db)
	uRepo := repository.GetUserRepository(db)
	aRepo := repository.GetAuditLogRepository(db)

	cUsecase := usecase.GetCategoryUsecase(aRepo, cRepo)
	pUsecase := usecase.GetProductUsecase(aRepo, pRepo, cRepo, tRepo)
	mUsecase := usecase.GetMemberUsecase(aRepo, mRepo, aesKey, bindexKey)
	oUsecase := usecase.GetOrderUsecase(aRepo, oRepo, mRepo, pRepo)
	tUsecase := usecase.GetTaxUsecase(aRepo, tRepo)
	uUsecase := usecase.GetUserUsecase(aRepo, uRepo, aesKey, bindexKey, rsaPrivateKey, rsaPublicKey)

	handler := rest.NewHandler(cUsecase, pUsecase, mUsecase, oUsecase, tUsecase, uUsecase)

	rest.LoadRoutes(e, handler, rdb)

	go func() {
		if err := e.Start(":3000"); err != nil && err != http.ErrServerClosed {
			slog.Error("SERVER: START SERVER FAILED", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("SERVER: GRACEFUL SHUTDOWN INITIATED...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		slog.Error("SERVER: SHUTDOWN ERROR", slog.String("error", err.Error()))
	}

	if err := rdb.Close(); err != nil {
		slog.Error("REDIS: CONNECTION CLOSE FAILED", slog.String("error", err.Error()))
	} else {
		slog.Info("REDIS: CONNECTION CLOSE SUCCESS")
	}

	sqlDB, err := db.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			slog.Error("POSTGRES: CONNECTION CLOSE FAILED", slog.String("error", err.Error()))
		} else {
			slog.Info("POSTGRES: CONNECTION CLOSE SUCCESS")
		}
	} else {
		slog.Error("POSTGRES: FAILED TO EXTRACT SQL DB INSTANCE", slog.String("error", err.Error()))
	}

	slog.Info("SERVER: SHUTDOWN COMPLETE")
}

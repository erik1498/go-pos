package main

import (
	"context"
	"go-pos/internal/database"
	"go-pos/internal/delivery/rest"
	"go-pos/internal/domain"
	"go-pos/internal/repository"
	"go-pos/internal/usecase"
	"go-pos/pkg/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
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
		log.Fatalf("REDIS: CONNECTION FAILED %s:%s. ERR: %v", host, port, err)
	}

	log.Println("REDIS: CONNECTED")
	return client
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("INFO: .env not found.")
	}

	aesKey := os.Getenv("AES_256_KEY")
	bindexKey := os.Getenv("BLIND_INDEX_KEY")

	if len(aesKey) != 32 {
		log.Fatal(domain.AlertAESNot32Character)
	}
	if bindexKey == "" {
		log.Fatal(domain.AlertBlindIndexEmpty)
	}

	rsaPrivateKey, err := utils.LoadRSAPrivateKey()
	if err != nil {
		panic(err)
	}

	rsaPublicKey, err := utils.LoadRSAPublicKey()
	if err != nil {
		panic(err)
	}

	dbAddress := os.Getenv("DATABASE_URL")
	if dbAddress == "" {
		dbAddress = "host=localhost user=postgres password=postgres dbname=go-pos sslmode=disable"
	}
	db := database.GetDB(dbAddress)
	rdb := initRedis()

	e := echo.New()

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
			e.Logger.Fatalf("SERVER: START SERVER FAILED, ERR: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("SERVER: GRACEFULL SHUTDOWN...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	if err := rdb.Close(); err != nil {
		log.Printf("REDIS: CONNECTION CLOSE FAILED, ERR: %v", err)
	} else {
		log.Println("REDIS: CONNECTION CLOSE SUCCESS")
	}

	log.Println("SERVER: SERVER SHUTDOWN SUCCESS")
}

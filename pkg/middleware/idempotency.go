package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"go-pos/internal/domain"
	"go-pos/pkg/response"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

const IdempotencyKeyCtx = "idempotency_key"

type IdempotencyResponse struct {
	StatusCode int             `json:"status_code"`
	Body       json.RawMessage `json:"body"`
}

type customResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *customResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Idempotency(rdb *redis.Client, idmpID string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			idemKey := c.Request().Header.Get("X-Idempotency-Key")
			if idemKey == "" {
				return response.ErrBadRequest(c, domain.ErrIdempotencyRequired.Error())
			}

			ctx := c.Request().Context()
			redisKey := "idmp:" + idmpID + ":" + idemKey

			cachedData, err := rdb.Get(ctx, redisKey).Result()
			if err != nil && err != redis.Nil {
				log.Printf("REDIS DISCONNECTED, ERR: %v", err)

				newCtx := context.WithValue(ctx, IdempotencyKeyCtx, idemKey)
				c.SetRequest(c.Request().WithContext(newCtx))

				return next(c)
			}

			if err == nil {
				if cachedData == "PROCESSING" {
					return response.ErrConflictRequest(c, domain.ErrRequestProcessed.Error())
				}

				var cachedResp IdempotencyResponse
				if err := json.Unmarshal([]byte(cachedData), &cachedResp); err == nil {
					return c.JSONBlob(cachedResp.StatusCode, cachedResp.Body)
				}
			}

			locked, err := rdb.SetNX(ctx, redisKey, "PROCESSING", 24*time.Hour).Result()
			if err != nil {
				log.Printf("REDIS SETNX FAILED, ERR: %v", err)
				newCtx := context.WithValue(ctx, IdempotencyKeyCtx, idemKey)
				c.SetRequest(c.Request().WithContext(newCtx))
				return next(c)
			}

			if !locked {
				return response.ErrConflictRequest(c, domain.ErrRequestProcessed.Error())
			}

			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &customResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer

			newCtx := context.WithValue(ctx, IdempotencyKeyCtx, idemKey)
			c.SetRequest(c.Request().WithContext(newCtx))

			err = next(c)

			if err == nil && c.Response().Status >= 200 && c.Response().Status < 300 {
				finalResp := IdempotencyResponse{
					StatusCode: c.Response().Status,
					Body:       resBody.Bytes(),
				}
				finalJSON, _ := json.Marshal(finalResp)
				go func(key string, data []byte) {
					_ = rdb.Set(context.Background(), key, data, 24*time.Hour).Err()
				}(redisKey, finalJSON)
			} else {
				go func(key string) {
					_ = rdb.Del(context.Background(), key).Err()
				}(redisKey)
			}

			return err
		}
	}
}

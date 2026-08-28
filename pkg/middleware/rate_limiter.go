package middleware

import (
	"fmt"
	"go-pos/internal/domain"
	"go-pos/pkg/response"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func RateLimiter(rdb *redis.Client, limit int, duration time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			clientIP := c.RealIP()
			redisKey := fmt.Sprintf("ratelimit:%s", clientIP)

			ctx := c.Request().Context()

			currentCount, err := rdb.Incr(ctx, redisKey).Result()
			if err != nil {
				return next(c)
			}

			if currentCount == 1 {
				rdb.Expire(ctx, redisKey, duration)
			}

			if int(currentCount) > limit {
				return response.ErrTooManyRequests(c, domain.ErrRateLimitExceeded.Error())
			}

			c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
			c.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", max(0, limit-int(currentCount))))

			return next(c)
		}
	}
}

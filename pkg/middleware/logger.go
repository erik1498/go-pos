package middleware

import (
	"log/slog"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

func InitSlog() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)
}

func SlogMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			reqID := res.Header().Get(echo.HeaderXRequestID)
			if reqID == "" {
				reqID = req.Header.Get(req.Header.Get(echo.HeaderXRequestID))
			}

			attrs := []slog.Attr{
				slog.String("request_id", reqID),
				slog.String("method", req.Method),
				slog.String("uri", req.RequestURI),
				slog.Int("status", res.Status),
				slog.String("latency", time.Since(start).String()),
				slog.String("ip", c.RealIP()),
				slog.String("user_agent", req.UserAgent()),
			}

			ctx := req.Context()
			if res.Status >= 500 {
				slog.LogAttrs(ctx, slog.LevelError, "SERVER ERROR", attrs...)
			} else if res.Status >= 400 {
				slog.LogAttrs(ctx, slog.LevelWarn, "CLIENT ERROR", attrs...)
			} else {
				slog.LogAttrs(ctx, slog.LevelInfo, "REQUEST PROCESSED", attrs...)
			}

			return nil
		}
	}
}

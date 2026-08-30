package middleware

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/labstack/echo/v4"
	"gopkg.in/natefinch/lumberjack.v2"
)

var sensitiveURIRegex = regexp.MustCompile(`(?i)(password|token|secret|key|pin)=([^&]+)`)

func InitSlog() {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0700); err != nil {
		panic("LOGGER: MAKE LOGS DIRECTORY FAIL, ERR: " + err.Error())
	}

	fileLogger := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, "app.log"),
		MaxSize:    10,
		MaxBackups: 7,
		MaxAge:     90,
		Compress:   true,
	}

	multiWriter := io.MultiWriter(os.Stdout, fileLogger)

	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.String(slog.TimeKey, a.Value.Time().UTC().Format(time.RFC3339))
			}
			return a
		},
	})

	slog.SetDefault(slog.New(handler))
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

			safeURI := sensitiveURIRegex.ReplaceAllString(req.RequestURI, "$1=***")

			attrs := []slog.Attr{
				slog.String("request_id", reqID),
				slog.String("method", req.Method),
				slog.String("uri", safeURI),
				slog.Int("status", res.Status),
				slog.String("latency", time.Since(start).String()),
				slog.String("ip", c.RealIP()),
				slog.String("user_agent", req.UserAgent()),
			}

			meta, metaValid := req.Context().Value(AuditMetaKey).(AuditMeta)
			if metaValid && meta.UserID != "" {
				attrs = append(attrs, slog.String("actor_id", meta.UserID))
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

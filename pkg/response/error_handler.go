package response

import (
	"errors"
	"go-pos/internal/domain"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var statusCode int
	var message string
	var logLevel = slog.LevelError

	if customErr, ok := errors.AsType[*domain.CustomError](err); ok {
		statusCode = customErr.HTTPCode
		message = customErr.Message

		if statusCode < 500 {
			logLevel = slog.LevelWarn
		}
	} else {
		if echoErr, ok := errors.AsType[*echo.HTTPError](err); ok {
			statusCode = echoErr.Code
			message = echoErr.Message.(string)
			if statusCode < 500 {
				logLevel = slog.LevelWarn
			}
		} else {
			statusCode = http.StatusInternalServerError
			message = domain.ErrInternalServer.Message
		}
	}

	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	if reqID == "" {
		reqID = c.Request().Header.Get(echo.HeaderXRequestID)
	}

	logAttrs := []slog.Attr{
		slog.String("request_id", reqID),
		slog.String("method", c.Request().Method),
		slog.String("path", c.Request().URL.Path),
		slog.String("error_trace", err.Error()),
	}

	if logLevel == slog.LevelWarn {
		slog.LogAttrs(c.Request().Context(), slog.LevelWarn, "Client Request Error", logAttrs...)
	} else {
		slog.LogAttrs(c.Request().Context(), slog.LevelError, "Internal Server Error", logAttrs...)
	}

	_ = c.JSON(statusCode, echo.Map{
		"meta": echo.Map{
			"code":    statusCode,
			"message": message,
		},
		"data": nil,
	})
}

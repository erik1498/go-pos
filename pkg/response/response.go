package response

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type SuccessPayload struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type ErrorPayload struct {
	Type     string      `json:"type"`
	Title    string      `json:"title"`
	Status   int         `json:"status"`
	Detail   string      `json:"detail"`
	Instance string      `json:"instance,omitempty"`
	Errors   interface{} `json:"errors,omitempty"`
}

type MetaPage struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalPages int   `json:"total_page"`
	TotalItems int64 `json:"total_items"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

func Success(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, SuccessPayload{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func NoContent(c echo.Context) error {
	return c.JSON(http.StatusNoContent, nil)
}

func SuccessWithMeta(c echo.Context, statusCode int, message string, data interface{}, meta MetaPage) error {
	return c.JSON(statusCode, SuccessPayload{
		Status:  "success",
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func BuildError(c echo.Context, statusCode int, title, detail string, errs interface{}) error {
	errType := strings.ToLower(strings.ReplaceAll(http.StatusText(statusCode), " ", "-"))

	payload := ErrorPayload{
		Type:     fmt.Sprintf("http://localhost:3000/errors/%s", errType),
		Title:    title,
		Status:   statusCode,
		Detail:   detail,
		Instance: c.Request().RequestURI,
		Errors:   errs,
	}

	c.Response().Header().Set(echo.HeaderContentType, "application/problem+json")
	return c.JSON(statusCode, payload)
}

func ErrBadRequest(c echo.Context, detail string) error {
	return BuildError(c, http.StatusBadRequest, "Bad Request", detail, nil)
}

func ErrConflictRequest(c echo.Context, detail string) error {
	return BuildError(c, http.StatusConflict, "Conflict", detail, nil)
}

func ErrValidation(c echo.Context, detail string, fieldErrors interface{}) error {
	return BuildError(c, http.StatusUnprocessableEntity, "Validation Failed", detail, fieldErrors)
}

func ErrNotFound(c echo.Context, detail string) error {
	return BuildError(c, http.StatusNotFound, "Not Found", detail, nil)
}

func ErrInternalServer(c echo.Context, detail string) error {
	return BuildError(c, http.StatusInternalServerError, "Internal Server Error", detail, nil)
}

package rest

import (
	"encoding/json"
	"errors"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) RegisterUser(c echo.Context) error {
	var req model.RegisterUser

	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return response.ErrBadRequest(c, domain.ErrBadRequest.Error())
	}

	user, err := h.uUsecase.Register(req)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyRegistered) {
			return response.ErrBadRequest(c, domain.ErrEmailAlreadyRegistered.Error())
		}
		return response.ErrInternalServer(c, err.Error())
	}

	return response.Success(c, http.StatusCreated, domain.SuccessCreateData, user)
}

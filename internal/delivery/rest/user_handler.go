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
	var req model.RegisterRequest

	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return response.ErrBadRequest(c, domain.ErrBadRequest.Error())
	}

	user, err := h.uUsecase.Register(req)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyRegistered) {
			return response.ErrBadRequest(c, domain.ErrEmailAlreadyRegistered.Error())
		}
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusCreated, domain.SuccessCreateData, user)
}

func (h *handler) LoginUser(c echo.Context) error {
	var req model.LoginRequest

	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return response.ErrBadRequest(c, domain.ErrBadRequest.Error())
	}

	userSession, err := h.uUsecase.Login(req)
	if err != nil {
		if errors.Is(err, domain.ErrUsernameOrPasswordInvalid) {
			return response.ErrBadRequest(c, domain.ErrUsernameOrPasswordInvalid.Error())
		}
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusOK, domain.SuccessLogin, userSession)
}

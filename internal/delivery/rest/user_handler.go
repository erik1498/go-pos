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

	ctx := c.Request().Context()

	user, err := h.uUsecase.Register(ctx, req)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyRegistered) {
			return response.ErrBadRequest(c, domain.ErrEmailAlreadyRegistered.Error())
		}
		if errors.Is(err, domain.ErrIdempotencyKeyDuplicate) {
			return response.ErrConflictRequest(c, domain.ErrIdempotencyKeyDuplicate.Error())
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

	ctx := c.Request().Context()

	userSession, err := h.uUsecase.Login(ctx, req)
	if err != nil {
		if errors.Is(err, domain.ErrUsernameOrPasswordInvalid) {
			return response.ErrBadRequest(c, domain.ErrUsernameOrPasswordInvalid.Error())
		}
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusOK, domain.SuccessLogin, userSession)
}

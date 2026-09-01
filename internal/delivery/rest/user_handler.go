package rest

import (
	"fmt"
	"go-pos/internal/domain"
	"go-pos/pkg/response"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserSessionResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func toUserSessionResponse(u domain.UserSession) UserSessionResponse {
	return UserSessionResponse{
		AccessToken:  u.AccessToken,
		RefreshToken: u.RefreshToken,
	}
}

func (h *handler) RegisterUser(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("[delivery][rest][user_handler][RegisterUser] invalid body: %w", domain.ErrBadRequest)
	}

	if err := c.Validate(&req); err != nil {
		return fmt.Errorf("[delivery][rest][user_handler][RegisterUser] validation error: %w", domain.ErrBadRequest)
	}

	ctx := c.Request().Context()

	param := domain.RegisterUserParam{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := h.uUsecase.Register(ctx, param)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusCreated, domain.SuccessCreateData, user)
}

func (h *handler) LoginUser(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("[delivery][rest][user_handler][RegisterUser] invalid body: %w", domain.ErrBadRequest)
	}

	if err := c.Validate(&req); err != nil {
		return fmt.Errorf("[delivery][rest][user_handler][RegisterUser] validation error: %w", domain.ErrBadRequest)
	}

	ctx := c.Request().Context()

	param := domain.LoginUserParam{
		Username: req.Username,
		Password: req.Password,
	}

	userSession, err := h.uUsecase.Login(ctx, param)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, domain.SuccessLogin, userSession)
}

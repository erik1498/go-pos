package domain

import (
	"context"
	"errors"
	"go-pos/internal/model"
)

var (
	ErrUsernameOrPasswordInvalid = errors.New("username or password invalid")
	ErrUnauthorized              = errors.New("unauthorized")
	ErrSessionExpired            = errors.New("user session invalid or expired")
	ErrForbidden                 = errors.New("access forbidden")
	ErrIdempotencyRequired       = errors.New("header x-idempotency-key required")
	ErrRequestProcessed          = errors.New("request on process, pelase wait")

	SuccessLogin = "login success"
)

type UserRepository interface {
	Register(ctx context.Context, user model.User) (model.User, error)
	GetByUsername(ctx context.Context, username string) (model.User, error)
}

type UserUsecase interface {
	Register(ctx context.Context, req model.RegisterRequest) (model.User, error)
	Login(ctx context.Context, req model.LoginRequest) (model.UserSession, error)
	CheckSession(ctx context.Context, userSession model.UserSession) (string, string, error)
}

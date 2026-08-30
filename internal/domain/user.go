package domain

import (
	"context"
	"go-pos/internal/model"
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

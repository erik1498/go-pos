package domain

import (
	"errors"
	"go-pos/internal/model"
)

var (
	ErrUsernameOrPasswordInvalid = errors.New("username or password invalid")
	ErrUnauthorized              = errors.New("unauthorized")
	ErrSessionExpired            = errors.New("user session invalid or expired")

	SuccessLogin = "login success"
)

type UserRepository interface {
	Register(user model.User) (model.User, error)
	GetByUsername(username string) (model.User, error)
}

type UserUsecase interface {
	Register(req model.RegisterRequest) (model.User, error)
	Login(req model.LoginRequest) (model.UserSession, error)
	CheckSession(userSession model.UserSession) (string, error)
}

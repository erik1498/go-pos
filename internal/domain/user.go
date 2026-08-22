package domain

import "go-pos/internal/model"

var (
	ErrUsernameOrPasswordInvalid = "username or password invalid"
)

type UserRepository interface {
	Register(user model.User) (model.User, error)
}

type UserUsecase interface {
	Register(req model.RegisterUser) (model.User, error)
}

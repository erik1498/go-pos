package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleCashier UserRole = "CASHIER"
	UserRoleAdmin   UserRole = "ADMIN"
	AuthContextKey  string   = "user_auth"
)

type User struct {
	ID             uuid.UUID
	Username       string
	EmailEncrypted []byte
	EmailBIndex    string
	Email          string
	Hash           string
	Role           UserRole
	IdempotencyKey string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type RegisterUserParam struct {
	Username string
	Email    string
	Password string
}

type LoginUserParam struct {
	Username string
	Password string
}

type UserSession struct {
	AccessToken  string
	RefreshToken string
}

type UserRepository interface {
	Register(ctx context.Context, user User) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
}

type UserUsecase interface {
	Register(ctx context.Context, req RegisterUserParam) (User, error)
	Login(ctx context.Context, req LoginUserParam) (UserSession, error)
	CheckSession(ctx context.Context, userSession UserSession) (string, string, error)
}

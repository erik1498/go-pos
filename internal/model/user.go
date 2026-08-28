package model

import (
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
	ID             uuid.UUID `gorm:"primaryKey;type:uuid;uniqueIndex;not null" json:"-"`
	Username       string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	EmailEncrypted []byte    `gorm:"type:bytea;not null;" json:"-"`
	EmailBIndex    string    `gorm:"type:varchar(64);uniqueIndex:idx_email_b_index;not null" json:"-"`
	Email          string    `gorm:"-" json:"email,omitempty"`
	Hash           string    `gorm:"type:varchar(255);not null;" json:"-"`
	Role           UserRole  `gorm:"type:varchar(20);not null;default:'CASHIER'" json:"role"`
	IdempotencyKey string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"idempotency_key"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type UserSession struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

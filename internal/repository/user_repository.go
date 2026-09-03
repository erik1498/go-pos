package repository

import (
	"context"
	"errors"
	"fmt"
	"go-pos/internal/domain"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type UserRole string

type UserDAO struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid;uniqueIndex;not null"`
	Username       string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	EmailEncrypted []byte    `gorm:"type:bytea;not null;"`
	EmailBIndex    string    `gorm:"type:varchar(64);uniqueIndex:idx_email_b_index;not null"`
	Email          string    `gorm:"-" json:"email,omitempty"`
	Hash           string    `gorm:"type:varchar(255);not null;"`
	Role           UserRole  `gorm:"type:varchar(20);not null;default:'CASHIER'"`
	IdempotencyKey string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (UserDAO) TableName() string {
	return "users"
}

func (dao *UserDAO) ToDomain() domain.User {
	return domain.User{
		ID:             dao.ID,
		Username:       dao.Username,
		Email:          dao.Email,
		EmailEncrypted: dao.EmailEncrypted,
		EmailBIndex:    dao.EmailBIndex,
		Hash:           dao.Hash,
		Role:           domain.UserRole(dao.Role),
		IdempotencyKey: dao.IdempotencyKey,
		CreatedAt:      dao.CreatedAt,
		UpdatedAt:      dao.UpdatedAt,
	}
}

func FromDomainUser(u domain.User) UserDAO {
	dao := UserDAO{
		ID:             u.ID,
		Username:       u.Username,
		EmailEncrypted: u.EmailEncrypted,
		EmailBIndex:    u.EmailBIndex,
		Email:          u.Email,
		Hash:           u.Hash,
		Role:           UserRole(u.Role),
		IdempotencyKey: u.IdempotencyKey,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
	return dao
}

type userRepository struct {
	db *gorm.DB
}

func GetUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (uRepo *userRepository) Register(ctx context.Context, user domain.User) (domain.User, error) {
	dao := FromDomainUser(user)

	if err := uRepo.db.WithContext(ctx).Create(&dao).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "idx_email_b_index" {
				return domain.User{}, fmt.Errorf("[repository][user_repository][Register] err idempotency key: %w", domain.ErrEmailAlreadyRegistered)
			}
		}
		return domain.User{}, fmt.Errorf("[repository][user_repository][Register] db query failed: %w", err)
	}
	return user, nil
}

func (uRepo *userRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	var user UserDAO
	if err := uRepo.db.WithContext(ctx).Where(&UserDAO{Username: username}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, fmt.Errorf("[repository][user_repository][GetByUsername] err not found: %w", domain.ErrUsernameOrPasswordInvalid)
		}
		return domain.User{}, fmt.Errorf("[repository][user_repository][GetByUsername] db query failed: %w", err)
	}
	return user.ToDomain(), nil
}

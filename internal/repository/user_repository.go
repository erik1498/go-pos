package repository

import (
	"errors"
	"go-pos/internal/domain"
	"go-pos/internal/model"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func GetUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (uRepo *userRepository) Register(user model.User) (model.User, error) {
	if err := uRepo.db.Create(&user).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "idx_email_b_index" {
				return model.User{}, domain.ErrEmailAlreadyRegistered
			}
		}
		return model.User{}, err
	}
	return user, nil
}

package usecase

import (
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/utils"

	"github.com/google/uuid"
)

type userUsecase struct {
	uRepo     domain.UserRepository
	aesKey    []byte
	bindexKey string
}

func GetUserUsecase(uRepo domain.UserRepository, aesKey, bindexKey string) domain.UserUsecase {
	return &userUsecase{
		uRepo:     uRepo,
		aesKey:    []byte(aesKey),
		bindexKey: bindexKey,
	}
}

func (uUsecase *userUsecase) Register(req model.RegisterUser) (model.User, error) {
	emailEncrypted, err := utils.EncryptAES(req.Email, uUsecase.aesKey)
	if err != nil {
		return model.User{}, domain.ErrEncryptEmail
	}

	passwordHash, err := utils.HashPassword(req.Password, uUsecase.aesKey)
	if err != nil {
		return model.User{}, err
	}

	user := model.User{
		ID:             uuid.Must(uuid.NewV7()),
		Username:       req.Username,
		EmailEncrypted: emailEncrypted,
		EmailBIndex:    utils.GenerateBlindedIndex(req.Email, uUsecase.bindexKey),
		Role:           model.UserRoleCashier,
		Email:          req.Email,
		Hash:           passwordHash,
	}

	user, err = uUsecase.uRepo.Register(user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

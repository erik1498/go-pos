package usecase

import (
	"crypto/rsa"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/utils"

	"github.com/google/uuid"
)

type userUsecase struct {
	uRepo     domain.UserRepository
	aesKey    []byte
	rsa256key *rsa.PrivateKey
	bindexKey string
}

func GetUserUsecase(uRepo domain.UserRepository, aesKey, bindexKey string, rsa256Key *rsa.PrivateKey) domain.UserUsecase {
	return &userUsecase{
		uRepo:     uRepo,
		aesKey:    []byte(aesKey),
		rsa256key: rsa256Key,
		bindexKey: bindexKey,
	}
}

func (uUsecase *userUsecase) Register(req model.RegisterRequest) (model.User, error) {
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

func (uUsecase *userUsecase) Login(req model.LoginRequest) (model.UserSession, error) {
	user, err := uUsecase.uRepo.GetByUsername(req.Username)
	if err != nil {
		return model.UserSession{}, err
	}

	if user.Username != req.Username {
		return model.UserSession{}, domain.ErrUsernameOrPasswordInvalid
	}

	verifyPassword, err := utils.VerifyPassword(req.Password, user.Hash, uUsecase.aesKey)
	if err != nil {
		return model.UserSession{}, domain.ErrUsernameOrPasswordInvalid
	}
	if !verifyPassword {
		return model.UserSession{}, domain.ErrUsernameOrPasswordInvalid
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, string(user.Role), uUsecase.rsa256key)
	if err != nil {
		return model.UserSession{}, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, uUsecase.rsa256key)
	if err != nil {
		return model.UserSession{}, err
	}

	return model.UserSession{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

package usecase

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/middleware"
	"go-pos/pkg/utils"
	"log"

	"github.com/google/uuid"
)

type userUsecase struct {
	aRepo            domain.AuditLogRepository
	uRepo            domain.UserRepository
	aesKey           []byte
	rsa256PrivateKey *rsa.PrivateKey
	rsa256PublicKey  *rsa.PublicKey
	bindexKey        string
}

func GetUserUsecase(aRepo domain.AuditLogRepository, uRepo domain.UserRepository, aesKey, bindexKey string, rsa256PrivateKey *rsa.PrivateKey, rsa256PublicKey *rsa.PublicKey) domain.UserUsecase {
	return &userUsecase{
		aRepo:            aRepo,
		uRepo:            uRepo,
		aesKey:           []byte(aesKey),
		rsa256PrivateKey: rsa256PrivateKey,
		rsa256PublicKey:  rsa256PublicKey,
		bindexKey:        bindexKey,
	}
}

func (uUsecase *userUsecase) Register(ctx context.Context, req model.RegisterRequest) (model.User, error) {
	meta, metaValid := ctx.Value(middleware.AuditMetaKey).(middleware.AuditMeta)
	var actorID uuid.UUID
	if metaValid {
		actorID = uuid.MustParse(meta.UserID)
	}

	idemKey, ok := ctx.Value(middleware.IdempotencyKeyCtx).(string)
	if !ok && idemKey == "" {
		return model.User{}, domain.ErrIdempotencyKeyDuplicate
	}

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
		IdempotencyKey: idemKey,
	}

	user, err = uUsecase.uRepo.Register(ctx, user)
	if err != nil {
		return model.User{}, err
	}

	if metaValid {
		newValueJSON, _ := json.Marshal(user)

		auditLog := model.AuditLog{
			ActorID:   actorID,
			ActorRole: meta.Role,
			Action:    "CREATE",
			Entity:    "users",
			EntityID:  user.ID.String(),
			OldValues: "null",
			NewValues: string(newValueJSON),
			IPAddress: meta.IPAddress,
			UserAgent: meta.UserAgent,
		}

		go func(logData model.AuditLog) {
			if err := uUsecase.aRepo.Create(context.Background(), logData); err != nil {
				log.Printf("AUDIT LOG: RECORD AUDIT LOG FAILED, ERR : %v", err)
			}
		}(auditLog)
	}

	return user, nil
}

func (uUsecase *userUsecase) Login(ctx context.Context, req model.LoginRequest) (model.UserSession, error) {
	user, err := uUsecase.uRepo.GetByUsername(ctx, req.Username)
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

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, string(user.Role), uUsecase.rsa256PrivateKey)
	if err != nil {
		return model.UserSession{}, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, uUsecase.rsa256PrivateKey)
	if err != nil {
		return model.UserSession{}, err
	}

	return model.UserSession{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uUsecase *userUsecase) CheckSession(ctx context.Context, userSession model.UserSession) (string, string, error) {
	userID, role, err := utils.GetClaimsFromAccessToken(userSession.AccessToken, uUsecase.rsa256PublicKey)
	if err != nil {
		return "", "", err
	}
	return userID, role, nil
}

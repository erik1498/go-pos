package usecase

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"go-pos/internal/domain"
	"go-pos/pkg/middleware"
	"go-pos/pkg/utils"
	"log/slog"

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

func (uUsecase *userUsecase) Register(ctx context.Context, req domain.RegisterUserParam) (domain.User, error) {
	meta, metaValid := ctx.Value(middleware.AuditMetaKey).(middleware.AuditMeta)
	var actorID uuid.UUID
	if metaValid {
		actorID = uuid.MustParse(meta.UserID)
	}

	idemKey, ok := ctx.Value(middleware.IdempotencyKeyCtx).(string)
	if !ok && idemKey == "" {
		return domain.User{}, domain.ErrIdempotencyKeyDuplicate
	}

	emailEncrypted, err := utils.EncryptAES(req.Email, uUsecase.aesKey)
	if err != nil {
		return domain.User{}, err
	}

	passwordHash, err := utils.HashPassword(req.Password, uUsecase.aesKey)
	if err != nil {
		return domain.User{}, err
	}

	user, err := uUsecase.uRepo.Register(ctx, domain.User{
		ID:             uuid.Must(uuid.NewV7()),
		Username:       req.Username,
		Email:          req.Email,
		EmailEncrypted: emailEncrypted,
		EmailBIndex:    utils.GenerateBlindedIndex(req.Email, uUsecase.bindexKey),
		Hash:           passwordHash,
		Role:           domain.UserRoleCashier,
		IdempotencyKey: idemKey,
	})

	if err != nil {
		return domain.User{}, err
	}

	if metaValid {
		newValuesJSON, errMarshal := json.Marshal(user)

		if errMarshal != nil {
			slog.Warn("[usecase][category_usecase][Register] failed to marshal new values",
				slog.String("error_trace", errMarshal.Error()),
			)
			newValuesJSON = []byte("{}")
		}

		auditLog := domain.AuditLog{
			ActorID:   actorID,
			ActorRole: meta.Role,
			Action:    "CREATE",
			Entity:    "users",
			EntityID:  user.ID.String(),
			OldValues: "null",
			NewValues: string(newValuesJSON),
			IPAddress: meta.IPAddress,
			UserAgent: meta.UserAgent,
		}

		go func(logData domain.AuditLog) {
			if err := uUsecase.aRepo.Create(context.Background(), logData); err != nil {
				slog.Error("[usecase][category_usecase][Register] failed to record audit log",
					slog.String("error_trace", err.Error()),
					slog.String("entity_id", logData.EntityID),
					slog.String("actor_id", logData.ActorID.String()),
				)
			}
		}(auditLog)
	}

	return user, nil
}

func (uUsecase *userUsecase) Login(ctx context.Context, req domain.LoginUserParam) (domain.UserSession, error) {
	user, err := uUsecase.uRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return domain.UserSession{}, err
	}

	if user.Username != req.Username {
		return domain.UserSession{}, domain.ErrUsernameOrPasswordInvalid
	}

	verifyPassword, err := utils.VerifyPassword(req.Password, user.Hash, uUsecase.aesKey)
	if err != nil {
		return domain.UserSession{}, domain.ErrUsernameOrPasswordInvalid
	}
	if !verifyPassword {
		return domain.UserSession{}, domain.ErrUsernameOrPasswordInvalid
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, string(user.Role), uUsecase.rsa256PrivateKey)
	if err != nil {
		return domain.UserSession{}, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, uUsecase.rsa256PrivateKey)
	if err != nil {
		return domain.UserSession{}, err
	}

	return domain.UserSession{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uUsecase *userUsecase) CheckSession(ctx context.Context, userSession domain.UserSession) (string, string, error) {
	userID, role, err := utils.GetClaimsFromAccessToken(userSession.AccessToken, uUsecase.rsa256PublicKey)
	if err != nil {
		return "", "", err
	}
	return userID, role, nil
}

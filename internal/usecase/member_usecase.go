package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-pos/internal/domain"
	"go-pos/pkg/middleware"
	"go-pos/pkg/utils"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

type memberUsecase struct {
	aRepo     domain.AuditLogRepository
	mRepo     domain.MemberRepository
	aesKey    []byte
	bindexKey string
}

func GetMemberUsecase(aRepo domain.AuditLogRepository, mRepo domain.MemberRepository, aesKey string, bindexKey string) domain.MemberUsecase {
	return &memberUsecase{
		aRepo:     aRepo,
		mRepo:     mRepo,
		aesKey:    []byte(aesKey),
		bindexKey: bindexKey,
	}
}

func (mUsecase *memberUsecase) GetAll(ctx context.Context, opts domain.QueryOptions) ([]domain.Member, int64, error) {
	allowedSorts := map[string]bool{
		"member_code": true,
		"created_at":  true,
	}

	allowedFields := map[string]bool{
		"member_code": true,
		"phone":       true,
		"email":       true,
	}

	cleanOpts := utils.SanitizeQuery(opts, allowedFields, allowedSorts, "created_at desc")

	for i, f := range cleanOpts.Filters {
		if f.Field == "phone" {
			cleanOpts.Filters[i].Field = "phone_bindex"
			cleanOpts.Filters[i].Value = utils.GenerateBlindedIndex(f.Value.(string), mUsecase.bindexKey)
			cleanOpts.Filters[i].Operator = "="
		}
		if f.Field == "email" {
			cleanOpts.Filters[i].Field = "email_bindex"
			cleanOpts.Filters[i].Value = utils.GenerateBlindedIndex(f.Value.(string), mUsecase.bindexKey)
			cleanOpts.Filters[i].Operator = "="
		}
	}

	encryptedMembers, totalItems, err := mUsecase.mRepo.GetAll(ctx, cleanOpts)
	if err != nil {
		return nil, 0, err
	}

	var decryptedMembers []domain.Member
	for _, m := range encryptedMembers {
		m.Name, err = utils.DecryptAES(m.NameEncrypted, mUsecase.aesKey)
		if err != nil {
			return nil, 0, err
		}
		m.Phone, err = utils.DecryptAES(m.PhoneEncrypted, mUsecase.aesKey)
		if err != nil {
			return nil, 0, err
		}
		m.Email, err = utils.DecryptAES(m.EmailEncrypted, mUsecase.aesKey)
		if err != nil {
			return nil, 0, err
		}
		decryptedMembers = append(decryptedMembers, m)
	}

	return decryptedMembers, totalItems, nil
}

func (mUsecase *memberUsecase) Create(ctx context.Context, req domain.CreateMemberParam) (domain.Member, error) {
	meta, metaValid := ctx.Value(middleware.AuditMetaKey).(middleware.AuditMeta)
	var actorID uuid.UUID
	if metaValid {
		actorID = uuid.MustParse(meta.UserID)
	}

	idemKey, ok := ctx.Value(middleware.IdempotencyKeyCtx).(string)
	if !ok && idemKey == "" {
		return domain.Member{}, domain.ErrIdempotencyKeyDuplicate
	}

	req.Phone = strings.TrimSpace(req.Phone)
	if req.Phone == "" {
		return domain.Member{}, errors.New("")
	}

	ID := uuid.Must(uuid.NewV7())
	memberCode := fmt.Sprintf("MBR-%s-%s", time.Now().Format("060102150405"), strings.ToUpper(ID.String()[:4]))

	nameEnc, err := utils.EncryptAES(req.Phone, mUsecase.aesKey)
	if err != nil {
		return domain.Member{}, err
	}

	phoneEnc, err := utils.EncryptAES(req.Phone, mUsecase.aesKey)
	if err != nil {
		return domain.Member{}, err
	}

	var emailEnc []byte
	if req.Email != "" {
		emailEnc, err = utils.EncryptAES(req.Email, mUsecase.aesKey)
		if err != nil {
			return domain.Member{}, err
		}
	}

	phoneBIndex := utils.GenerateBlindedIndex(req.Phone, mUsecase.bindexKey)

	var emailBIndex *string
	if req.Email == "" {
		hashEmail := utils.GenerateBlindedIndex(strings.ToLower(req.Email), mUsecase.bindexKey)
		emailBIndex = &hashEmail
	}

	newMember := domain.Member{
		ID:             ID,
		MemberCode:     memberCode,
		NameEncrypted:  nameEnc,
		PhoneEncrypted: phoneEnc,
		EmailEncrypted: emailEnc,
		PhoneBIndex:    phoneBIndex,
		EmailBIndex:    emailBIndex,
		Name:           req.Name,
		Phone:          req.Phone,
		Email:          req.Email,
		Points:         0,
		IdempotencyKey: idemKey,
		CreatedBy:      actorID,
		UpdatedBy:      actorID,
	}

	member, err := mUsecase.mRepo.Create(ctx, newMember)
	if err != nil {
		return domain.Member{}, err
	}

	if metaValid {
		newValueJSON, _ := json.Marshal(member)

		auditLog := domain.AuditLog{
			ActorID:   actorID,
			ActorRole: meta.Role,
			Action:    "CREATE",
			Entity:    "members",
			EntityID:  member.ID.String(),
			OldValues: "null",
			NewValues: string(newValueJSON),
			IPAddress: meta.IPAddress,
			UserAgent: meta.UserAgent,
		}

		go func(logData domain.AuditLog) {
			if err := mUsecase.aRepo.Create(context.Background(), logData); err != nil {
				log.Printf("AUDIT LOG: RECORD AUDIT LOG FAILED, ERR : %v", err)
			}
		}(auditLog)
	}

	return member, nil
}

func (mUsecase *memberUsecase) GetByID(ctx context.Context, id uuid.UUID) (domain.Member, error) {
	member, err := mUsecase.mRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Member{}, err
	}

	member.Name, err = utils.DecryptAES(member.NameEncrypted, mUsecase.aesKey)
	if err != nil {
		return domain.Member{}, err
	}
	member.Email, err = utils.DecryptAES(member.EmailEncrypted, mUsecase.aesKey)
	if err != nil {
		return domain.Member{}, err
	}
	member.Phone, err = utils.DecryptAES(member.PhoneEncrypted, mUsecase.aesKey)
	if err != nil {
		return domain.Member{}, err
	}

	return member, nil
}

func (mUsecase *memberUsecase) UpdateByID(ctx context.Context, id uuid.UUID, req domain.UpdateMemberParam) (domain.Member, error) {
	meta, metaValid := ctx.Value(middleware.AuditMetaKey).(middleware.AuditMeta)
	var actorID uuid.UUID
	if metaValid {
		actorID = uuid.MustParse(meta.UserID)
	}

	oldMember, err := mUsecase.mRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Member{}, err
	}

	nameEncrypted, err := utils.EncryptAES(req.Name, mUsecase.aesKey)

	if err != nil {
		return domain.Member{}, err
	}
	phoneEncrypted, err := utils.EncryptAES(req.Phone, mUsecase.aesKey)
	if err != nil {
		return domain.Member{}, err
	}
	emailEncrypted, err := utils.EncryptAES(req.Email, mUsecase.aesKey)
	if err != nil {
		return domain.Member{}, err
	}

	phoneBindex := utils.GenerateBlindedIndex(req.Phone, mUsecase.bindexKey)
	var emailBIndex *string
	if req.Email != "" {
		hashEmail := utils.GenerateBlindedIndex(req.Email, mUsecase.bindexKey)
		emailBIndex = &hashEmail
	}

	member := domain.Member{
		NameEncrypted:  nameEncrypted,
		PhoneEncrypted: phoneEncrypted,
		EmailEncrypted: emailEncrypted,
		EmailBIndex:    emailBIndex,
		PhoneBIndex:    phoneBindex,
		Name:           req.Name,
		Phone:          req.Phone,
		Email:          req.Email,
		UpdatedBy:      actorID,
	}

	updatedMember, err := mUsecase.mRepo.UpdateByID(ctx, id, member)

	if metaValid {
		oldValuesJSON, _ := json.Marshal(oldMember)
		newValuesJSON, _ := json.Marshal(updatedMember)

		auditLog := domain.AuditLog{
			ActorID:   actorID,
			ActorRole: meta.Role,
			Action:    "UPDATE",
			Entity:    "members",
			EntityID:  id.String(),
			OldValues: string(oldValuesJSON),
			NewValues: string(newValuesJSON),
			IPAddress: meta.IPAddress,
			UserAgent: meta.UserAgent,
		}

		go func(logData domain.AuditLog) {
			if err := mUsecase.aRepo.Create(context.Background(), logData); err != nil {
				log.Printf("AUDIT LOG: RECORD AUDIT LOG FAILED, ERR : %v", err)
			}
		}(auditLog)
	}

	return updatedMember, nil
}

func (mUsecase *memberUsecase) DeleteByID(ctx context.Context, id uuid.UUID) error {
	meta, metaValid := ctx.Value(middleware.AuditMetaKey).(middleware.AuditMeta)
	var actorID uuid.UUID
	if metaValid {
		actorID = uuid.MustParse(meta.UserID)
	}

	err := mUsecase.mRepo.DeleteByID(ctx, id, actorID)
	if err != nil {
		return err
	}

	if metaValid {
		auditLog := domain.AuditLog{
			ActorID:   actorID,
			ActorRole: meta.Role,
			Action:    "DELETE",
			Entity:    "members",
			EntityID:  id.String(),
			OldValues: "{}",
			NewValues: "null",
			IPAddress: meta.IPAddress,
			UserAgent: meta.UserAgent,
		}
		go func(logData domain.AuditLog) {
			mUsecase.aRepo.Create(context.Background(), logData)
		}(auditLog)
	}

	return nil
}

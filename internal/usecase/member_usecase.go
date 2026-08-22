package usecase

import (
	"fmt"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/utils"
	"strings"
	"time"

	"github.com/google/uuid"
)

type memberUsecase struct {
	mRepo     domain.MemberRepository
	aesKey    []byte
	bindexKey string
}

func GetMemberUsecase(mRepo domain.MemberRepository, aesKey string, bindexKey string) domain.MemberUsecase {
	return &memberUsecase{
		mRepo:     mRepo,
		aesKey:    []byte(aesKey),
		bindexKey: bindexKey,
	}
}

func (mUsecase *memberUsecase) GetAll(opts domain.QueryOptions) ([]model.Member, int64, error) {
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

	encryptedMembers, totalItems, err := mUsecase.mRepo.GetAll(cleanOpts)
	if err != nil {
		return nil, 0, err
	}

	var decryptedMembers []model.Member
	for _, m := range encryptedMembers {
		m.Name, err = utils.DecryptAES(m.NameEncrypted, mUsecase.aesKey)
		if err != nil {
			return nil, 0, domain.ErrDecryptName
		}
		m.Phone, err = utils.DecryptAES(m.PhoneEncrypted, mUsecase.aesKey)
		if err != nil {
			return nil, 0, domain.ErrDecryptPhone
		}
		m.Email, err = utils.DecryptAES(m.EmailEncrypted, mUsecase.aesKey)
		if err != nil {
			return nil, 0, domain.ErrDecryptEmail
		}
		decryptedMembers = append(decryptedMembers, m)
	}

	return decryptedMembers, totalItems, nil
}

func (mUsecase *memberUsecase) Create(req model.MemberRequest) (model.Member, error) {
	req.Phone = strings.TrimSpace(req.Phone)
	if req.Phone == "" {
		return model.Member{}, domain.ErrPhoneNumberRequired
	}

	ID := uuid.Must(uuid.NewV7())
	memberCode := fmt.Sprintf("MBR-%s-%s", time.Now().Format("060102150405"), strings.ToUpper(ID.String()[:4]))

	nameEnc, err := utils.EncryptAES(req.Phone, mUsecase.aesKey)
	if err != nil {
		return model.Member{}, domain.ErrEncryptName
	}

	phoneEnc, err := utils.EncryptAES(req.Phone, mUsecase.aesKey)
	if err != nil {
		return model.Member{}, domain.ErrEncryptPhone
	}

	var emailEnc []byte
	if req.Email != "" {
		emailEnc, err = utils.EncryptAES(req.Email, mUsecase.aesKey)
		if err != nil {
			return model.Member{}, domain.ErrEncryptEmail
		}
	}

	phoneBIndex := utils.GenerateBlindedIndex(req.Phone, mUsecase.bindexKey)

	var emailBIndex *string
	if req.Email == "" {
		hashEmail := utils.GenerateBlindedIndex(strings.ToLower(req.Email), mUsecase.bindexKey)
		emailBIndex = &hashEmail
	}

	newMember := model.Member{
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
	}

	member, err := mUsecase.mRepo.Create(newMember)
	if err != nil {
		return model.Member{}, err
	}

	return member, nil
}

func (mUsecase *memberUsecase) GetByID(id uuid.UUID) (model.Member, error) {
	member, err := mUsecase.mRepo.GetByID(id)
	if err != nil {
		return model.Member{}, err
	}

	member.Name, err = utils.DecryptAES(member.NameEncrypted, mUsecase.aesKey)
	if err != nil {
		return model.Member{}, domain.ErrDecryptName
	}
	member.Email, err = utils.DecryptAES(member.EmailEncrypted, mUsecase.aesKey)
	if err != nil {
		return model.Member{}, domain.ErrDecryptEmail
	}
	member.Phone, err = utils.DecryptAES(member.PhoneEncrypted, mUsecase.aesKey)
	if err != nil {
		return model.Member{}, domain.ErrDecryptPhone
	}

	return member, nil
}

func (mUsecase *memberUsecase) UpdateByID(req model.MemberRequest, id uuid.UUID) (model.Member, error) {
	nameEncrypted, err := utils.EncryptAES(req.Name, mUsecase.aesKey)
	if err != nil {
		return model.Member{}, domain.ErrEncryptName
	}
	phoneEncrypted, err := utils.EncryptAES(req.Phone, mUsecase.aesKey)
	if err != nil {
		return model.Member{}, domain.ErrEncryptName
	}
	emailEncrypted, err := utils.EncryptAES(req.Email, mUsecase.aesKey)
	if err != nil {
		return model.Member{}, domain.ErrEncryptName
	}

	phoneBindex := utils.GenerateBlindedIndex(req.Phone, mUsecase.bindexKey)
	var emailBIndex *string
	if req.Email != "" {
		hashEmail := utils.GenerateBlindedIndex(req.Email, mUsecase.bindexKey)
		emailBIndex = &hashEmail
	}

	member := model.Member{
		NameEncrypted:  nameEncrypted,
		PhoneEncrypted: phoneEncrypted,
		EmailEncrypted: emailEncrypted,
		EmailBIndex:    emailBIndex,
		PhoneBIndex:    phoneBindex,
		Name:           req.Name,
		Phone:          req.Phone,
		Email:          req.Email,
	}

	return mUsecase.mRepo.UpdateByID(member, id)
}

func (mUsecase *memberUsecase) DeleteByID(id uuid.UUID) error {
	return mUsecase.mRepo.DeleteByID(id)
}

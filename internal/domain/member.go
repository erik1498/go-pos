package domain

import (
	"errors"
	"go-pos/internal/model"
)

var (
	ErrPhoneNumberRequired    = errors.New("phone number is required")
	ErrEncryptName            = errors.New("name failed to encrypt")
	ErrEncryptPhone           = errors.New("phone failed to encrypt")
	ErrEncryptEmail           = errors.New("email failed to encrypt")
	ErrPhoneAlreadyRegistered = errors.New("phone already registered")
	ErrEmailAlreadyRegistered = errors.New("email already registered")

	ErrDecryptName  = errors.New("name failed to decrypt")
	ErrDecryptPhone = errors.New("phone failed to decrypt")
	ErrDecryptEmail = errors.New("email failed to decrypt")
)

type MemberRepository interface {
	GetAll(opts QueryOptions) ([]model.Member, int64, error)
	Create(member model.Member) (model.Member, error)
}

type MemberUsecase interface {
	GetAll(opts QueryOptions) ([]model.Member, int64, error)
	Create(member model.MemberRequest) (model.Member, error)
}

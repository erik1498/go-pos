package domain

import (
	"context"
	"errors"
	"go-pos/internal/model"

	"github.com/google/uuid"
)

var (
	ErrMemberNotFound = errors.New("member data not found")

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
	GetAll(ctx context.Context, opts QueryOptions) ([]model.Member, int64, error)
	Create(ctx context.Context, member model.Member) (model.Member, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Member, error)
	UpdateByID(ctx context.Context, id uuid.UUID, member model.Member) (model.Member, error)
	DeleteByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
}

type MemberUsecase interface {
	GetAll(ctx context.Context, opts QueryOptions) ([]model.Member, int64, error)
	Create(ctx context.Context, member model.MemberRequest) (model.Member, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Member, error)
	UpdateByID(ctx context.Context, id uuid.UUID, req model.MemberRequest) (model.Member, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

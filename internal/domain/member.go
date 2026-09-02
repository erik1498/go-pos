package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Member struct {
	ID             uuid.UUID
	MemberCode     string
	NameEncrypted  []byte
	PhoneEncrypted []byte
	EmailEncrypted []byte
	PhoneBIndex    string
	EmailBIndex    *string
	Name           string
	Phone          string
	Email          string
	Points         int
	IdempotencyKey string
	CreatedBy      uuid.UUID
	UpdatedBy      uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreateMemberParam struct {
	Name  string
	Phone string
	Email string
}

type UpdateMemberParam struct {
	Name  string
	Phone string
	Email string
}

type MemberRepository interface {
	GetAll(ctx context.Context, opts QueryOptions) ([]Member, int64, error)
	Create(ctx context.Context, member Member) (Member, error)
	GetByID(ctx context.Context, id uuid.UUID) (Member, error)
	UpdateByID(ctx context.Context, id uuid.UUID, member Member) (Member, error)
}

type MemberUsecase interface {
	GetAll(ctx context.Context, opts QueryOptions) ([]Member, int64, error)
	Create(ctx context.Context, member CreateMemberParam) (Member, error)
	GetByID(ctx context.Context, id uuid.UUID) (Member, error)
	UpdateByID(ctx context.Context, id uuid.UUID, req UpdateMemberParam) (Member, error)
}

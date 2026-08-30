package domain

import (
	"context"
	"go-pos/internal/model"

	"github.com/google/uuid"
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

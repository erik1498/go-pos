package repository

import (
	"context"
	"errors"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type memberRepository struct {
	db *gorm.DB
}

func GetMemberRepository(db *gorm.DB) domain.MemberRepository {
	return &memberRepository{
		db: db,
	}
}

func (mRepo *memberRepository) GetAll(ctx context.Context, opts domain.QueryOptions) ([]model.Member, int64, error) {
	var member []model.Member
	var totalItems int64

	dbQuery := mRepo.db.Model(&model.Member{})

	if opts.Search != "" {
		dbQuery = dbQuery.Where("member_code ILIKE ?", "%"+opts.Search+"%")
	}

	for _, f := range opts.Filters {
		queryStr := f.Field + " " + f.Operator + " ?"
		dbQuery.Where(queryStr, f.Value)
	}

	if err := dbQuery.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	err := dbQuery.Order(opts.Sort).Scopes(utils.PaginationScope(opts.Page, opts.Limit)).Find(&member).Error

	return member, totalItems, err
}

func (mRepo *memberRepository) Create(ctx context.Context, member model.Member) (model.Member, error) {
	if err := mRepo.db.Create(&member).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "idx_active_phone_bindex" {
				return model.Member{}, domain.ErrPhoneAlreadyRegistered
			}
			if pgErr.ConstraintName == "idx_active_email_bindex" {
				return model.Member{}, domain.ErrEmailAlreadyRegistered
			}
		}
		return model.Member{}, err
	}
	return member, nil
}

func (mRepo *memberRepository) GetByID(ctx context.Context, id uuid.UUID) (model.Member, error) {
	var member model.Member

	if err := mRepo.db.First(&member, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Member{}, domain.ErrMemberNotFound
		}
	}

	return member, nil
}

func (mRepo *memberRepository) UpdateByID(ctx context.Context, id uuid.UUID, member model.Member) (model.Member, error) {
	res := mRepo.db.Where(&model.Member{ID: id}).Clauses(clause.Returning{}).Updates(&member)

	if res.Error != nil {
		return model.Member{}, res.Error
	}

	if res.RowsAffected == 0 {
		return model.Member{}, domain.ErrMemberNotFound
	}
	return member, nil
}

func (mRepo *memberRepository) DeleteByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	err := mRepo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Member{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}

		res := tx.Delete(&model.Member{}, id)

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return domain.ErrMemberNotFound
		}

		return nil
	})

	return err
}

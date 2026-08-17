package repository

import (
	"errors"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/utils"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type memberRepository struct {
	db *gorm.DB
}

func GetMemberRepository(db *gorm.DB) domain.MemberRepository {
	return &memberRepository{
		db: db,
	}
}

func (mRepo *memberRepository) GetAll(opts domain.QueryOptions) ([]model.Member, int64, error) {
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

func (mRepo *memberRepository) Create(member model.Member) (model.Member, error) {
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

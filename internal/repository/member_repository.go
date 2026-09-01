package repository

import (
	"context"
	"errors"
	"fmt"
	"go-pos/internal/domain"
	"go-pos/pkg/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MemberDAO struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid();uniqueIndex;not null"`
	MemberCode     string    `gorm:"uniqueIndex;not null;type:varchar(50)"`
	NameEncrypted  []byte    `gorm:"type:bytea;not null"`
	PhoneEncrypted []byte    `gorm:"type:bytea;not null"`
	EmailEncrypted []byte    `gorm:"type:bytea;not null"`
	PhoneBIndex    string    `gorm:"type:varchar(64);uniqueIndex:idx_active_phone_bindex,where:deleted_at IS NULL;not null"`
	EmailBIndex    *string   `gorm:"type:varchar(64);uniqueIndex:idx_active_email_bindex,where:deleted_at IS NULL"`
	Name           string    `gorm:"-"`
	Phone          string    `gorm:"-"`
	Email          string    `gorm:"-"`
	Points         int       `gorm:"not null;default:0;"`
	IdempotencyKey string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	CreatedBy      uuid.UUID `gorm:"type:uuid;not null"`
	UpdatedBy      uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (MemberDAO) TableName() string {
	return "members"
}

func (dao *MemberDAO) ToDomain() domain.Member {
	return domain.Member{
		ID:             dao.ID,
		MemberCode:     dao.MemberCode,
		NameEncrypted:  dao.NameEncrypted,
		PhoneEncrypted: dao.PhoneEncrypted,
		EmailEncrypted: dao.EmailEncrypted,
		PhoneBIndex:    dao.PhoneBIndex,
		EmailBIndex:    dao.EmailBIndex,
		Name:           dao.Name,
		Phone:          dao.Phone,
		Email:          dao.Email,
		Points:         dao.Points,
		IdempotencyKey: dao.IdempotencyKey,
		CreatedBy:      dao.CreatedBy,
		UpdatedBy:      dao.UpdatedBy,
		CreatedAt:      dao.CreatedAt,
		UpdatedAt:      dao.UpdatedAt,
	}
}

func FromDomainMember(m domain.Member) MemberDAO {
	dao := MemberDAO{
		ID:             m.ID,
		MemberCode:     m.MemberCode,
		NameEncrypted:  m.NameEncrypted,
		PhoneEncrypted: m.PhoneEncrypted,
		EmailEncrypted: m.EmailEncrypted,
		PhoneBIndex:    m.PhoneBIndex,
		EmailBIndex:    m.EmailBIndex,
		Name:           m.Name,
		Phone:          m.Phone,
		Email:          m.Email,
		Points:         m.Points,
		IdempotencyKey: m.IdempotencyKey,
		CreatedBy:      m.CreatedBy,
		UpdatedBy:      m.UpdatedBy,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
	return dao
}

type memberRepository struct {
	db *gorm.DB
}

func GetMemberRepository(db *gorm.DB) domain.MemberRepository {
	return &memberRepository{
		db: db,
	}
}

func (mRepo *memberRepository) GetAll(ctx context.Context, opts domain.QueryOptions) ([]domain.Member, int64, error) {
	var daoList []MemberDAO
	var totalItems int64

	dbQuery := mRepo.db.Model(&MemberDAO{})

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

	if err := dbQuery.Order(opts.Sort).Scopes(utils.PaginationScope(opts.Page, opts.Limit)).Find(&daoList).Error; err != nil {
		return nil, 0, fmt.Errorf("[repository][member_repository][GetAll] db query failed: %w", err)
	}

	var members []domain.Member
	for _, m := range daoList {
		members = append(members, m.ToDomain())
	}

	return members, totalItems, nil
}

func (mRepo *memberRepository) Create(ctx context.Context, member domain.Member) (domain.Member, error) {
	dao := FromDomainMember(member)

	if err := mRepo.db.WithContext(ctx).Create(&dao).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "idx_categories_idempotency_key" {
				return domain.Member{}, fmt.Errorf("[repository][member_repository][Create] err idempotency key: %w", domain.ErrIdempotencyKeyDuplicate)
			}
			if pgErr.ConstraintName == "idx_active_phone_bindex" {
				return domain.Member{}, fmt.Errorf("[repository][member_repository][Create] err phone key: %w", domain.ErrPhoneAlreadyRegistered)
			}
			if pgErr.ConstraintName == "idx_active_email_bindex" {
				return domain.Member{}, fmt.Errorf("[repository][member_repository][Create] err email key: %w", domain.ErrEmailAlreadyRegistered)
			}
		}
		return domain.Member{}, fmt.Errorf("[repository][member_repository][Create] db query failed: %w", err)
	}
	return dao.ToDomain(), nil
}

func (mRepo *memberRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Member, error) {
	var member MemberDAO

	if err := mRepo.db.First(&member, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Member{}, fmt.Errorf("[repository][member_repository][GetByID] record not found: %w", domain.ErrMemberNotFound)
		}
		return domain.Member{}, fmt.Errorf("[repository][member_repository][GetByID] db query failed: %w", err)
	}

	return member.ToDomain(), nil
}

func (mRepo *memberRepository) UpdateByID(ctx context.Context, id uuid.UUID, member domain.Member) (domain.Member, error) {
	dao := FromDomainMember(member)

	res := mRepo.db.WithContext(ctx).Where("id = ?", id).Clauses(clause.Returning{}).Updates(&dao)

	if res.Error != nil {
		return domain.Member{}, fmt.Errorf("[repository][member_repository][UpdateByID] db query failed: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return domain.Member{}, fmt.Errorf("[repository][member_repository][UpdateByID] record not found: %w", domain.ErrMemberNotFound)
	}
	return dao.ToDomain(), nil
}

func (mRepo *memberRepository) DeleteByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	err := mRepo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&MemberDAO{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return fmt.Errorf("[repository][member_repository][DeleteMemberByID] db query failed: %w", err)
		}

		res := tx.Delete(&MemberDAO{}, id)

		if res.Error != nil {
			return fmt.Errorf("[repository][member_repository][DeleteMemberByID] db query failed: %w", res.Error)
		}

		if res.RowsAffected == 0 {
			return fmt.Errorf("[repository][member_repository][DeleteMemberByID] record not found: %w", domain.ErrCategoryNotFound)
		}

		return nil
	})

	return err
}

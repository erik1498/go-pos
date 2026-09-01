package repository

import (
	"context"
	"errors"
	"go-pos/internal/domain"
	"go-pos/pkg/utils"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TaxDAO struct {
	ID             uuid.UUID       `gorm:"primaryKey;type:uuid;default:gen_random_uuid();uniqueIndex;not null"`
	Name           string          `gorm:"type:varchar(50);uniqueIndex:idx_active_tax_name,where:deleted_at IS NULL;not null"`
	Rate           decimal.Decimal `gorm:"type:numeric(5,2);not null"`
	IsActive       bool            `gorm:"type:boolean;default:true"`
	IdempotencyKey string          `gorm:"type:varchar(100);uniqueIndex;not null"`
	CreatedBy      uuid.UUID       `gorm:"type:uuid;not null"`
	UpdatedBy      uuid.UUID       `gorm:"type:uuid;not null"`
	DeletedBy      *uuid.UUID      `gorm:"type:uuid"`
	CreatedAt      time.Time       `gorm:"autoCreateTime"`
	UpdatedAt      time.Time       `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt  `gorm:"index"`
}

func (TaxDAO) TableName() string {
	return "taxes"
}

func (dao *TaxDAO) ToDomain() domain.Tax {
	var deletedAt *time.Time
	if dao.DeletedAt.Valid {
		deletedAt = &dao.DeletedAt.Time
	}
	return domain.Tax{
		ID:             dao.ID,
		Name:           dao.Name,
		Rate:           dao.Rate,
		IsActive:       dao.IsActive,
		IdempotencyKey: dao.IdempotencyKey,
		CreatedBy:      dao.CreatedBy,
		UpdatedBy:      dao.UpdatedBy,
		DeletedBy:      dao.DeletedBy,
		CreatedAt:      dao.CreatedAt,
		UpdatedAt:      dao.UpdatedAt,
		DeletedAt:      deletedAt,
	}
}

func FromDomainTax(c domain.Tax) TaxDAO {
	dao := TaxDAO{
		ID:             c.ID,
		Name:           c.Name,
		Rate:           c.Rate,
		IsActive:       c.IsActive,
		IdempotencyKey: c.IdempotencyKey,
		CreatedBy:      c.CreatedBy,
		UpdatedBy:      c.UpdatedBy,
		DeletedBy:      c.DeletedBy,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
	if c.DeletedAt != nil {
		dao.DeletedAt = gorm.DeletedAt{Time: *c.DeletedAt, Valid: true}
	}
	return dao
}

type taxRepository struct {
	db *gorm.DB
}

func GetTaxRepository(db *gorm.DB) domain.TaxRepository {
	return &taxRepository{
		db: db,
	}
}

func (tRepo *taxRepository) GetAll(ctx context.Context, opts domain.QueryOptions) ([]domain.Tax, int64, error) {
	var taxes []domain.Tax
	var totalItems int64

	dbQuery := tRepo.db.Model(&domain.Tax{})
	if opts.Search != "" {
		dbQuery = dbQuery.Where("name ILIKE ?", "%"+opts.Search+"%")
	}

	for _, f := range opts.Filters {
		queryStr := f.Field + " " + f.Operator + " ?"
		dbQuery = dbQuery.Where(queryStr, f.Value)
	}

	if err := dbQuery.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	if err := dbQuery.Order(opts.Sort).Scopes(utils.PaginationScope(opts.Page, opts.Limit)).Find(&taxes).Error; err != nil {
		return nil, 0, err
	}

	return taxes, totalItems, nil
}

func (tRepo *taxRepository) Create(ctx context.Context, tax domain.Tax) (domain.Tax, error) {
	if err := tRepo.db.Create(&tax).Error; err != nil {
		return domain.Tax{}, err
	}
	return tax, nil
}

func (tRepo *taxRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Tax, error) {
	var tax domain.Tax

	if err := tRepo.db.First(&tax, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Tax{}, domain.ErrTaxNotFound
		}
		return domain.Tax{}, err
	}

	return tax, nil
}

func (tRepo *taxRepository) UpdateByID(ctx context.Context, id uuid.UUID, tax domain.Tax) (domain.Tax, error) {
	res := tRepo.db.Where(&domain.Tax{ID: id}).Clauses(clause.Returning{}).Updates(&tax)

	if res.Error != nil {
		return domain.Tax{}, res.Error
	}

	if res.RowsAffected == 0 {
		return domain.Tax{}, domain.ErrTaxNotFound
	}

	return tax, nil
}

func (tRepo *taxRepository) DeleteByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	err := tRepo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.Tax{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}

		res := tx.Delete(&domain.Tax{}, id)

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return domain.ErrTaxNotFound
		}

		return nil
	})

	return err
}

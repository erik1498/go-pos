package repository

import (
	"context"
	"errors"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type taxRepository struct {
	db *gorm.DB
}

func GetTaxRepository(db *gorm.DB) domain.TaxRepository {
	return &taxRepository{
		db: db,
	}
}

func (tRepo *taxRepository) GetAll(ctx context.Context, opts domain.QueryOptions) ([]model.Tax, int64, error) {
	var taxes []model.Tax
	var totalItems int64

	dbQuery := tRepo.db.Model(&model.Tax{})
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

func (tRepo *taxRepository) Create(ctx context.Context, tax model.Tax) (model.Tax, error) {
	if err := tRepo.db.Create(&tax).Error; err != nil {
		return model.Tax{}, err
	}
	return tax, nil
}

func (tRepo *taxRepository) GetByID(ctx context.Context, id uuid.UUID) (model.Tax, error) {
	var tax model.Tax

	if err := tRepo.db.First(&tax, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Tax{}, domain.ErrTaxNotFound
		}
		return model.Tax{}, err
	}

	return tax, nil
}

func (tRepo *taxRepository) UpdateByID(ctx context.Context, id uuid.UUID, tax model.Tax) (model.Tax, error) {
	res := tRepo.db.Where(&model.Tax{ID: id}).Clauses(clause.Returning{}).Updates(&tax)

	if res.Error != nil {
		return model.Tax{}, res.Error
	}

	if res.RowsAffected == 0 {
		return model.Tax{}, domain.ErrTaxNotFound
	}

	return tax, nil
}

func (tRepo *taxRepository) DeleteByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	err := tRepo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Tax{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}

		res := tx.Delete(&model.Tax{}, id)

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

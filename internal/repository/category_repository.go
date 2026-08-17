package repository

import (
	"errors"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type categoryRepo struct {
	db *gorm.DB
}

func GetCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepo{
		db: db,
	}
}

func (cRepo *categoryRepo) GetAll(opts domain.QueryOptions) ([]model.Category, int64, error) {
	var categoryList = []model.Category{}
	var totalItems int64

	dbQuery := cRepo.db.Model(&model.Category{})

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

	err := dbQuery.Order(opts.Sort).Scopes(utils.PaginationScope(opts.Page, opts.Limit)).Find(&categoryList).Error

	return categoryList, totalItems, err
}

func (cRepo *categoryRepo) Create(category model.Category) (model.Category, error) {
	if err := cRepo.db.Create(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Category{}, domain.CategoryErrNotFound
		}
		return model.Category{}, err
	}

	return category, nil
}

func (cRepo *categoryRepo) GetByID(id uuid.UUID) (model.Category, error) {
	var category model.Category

	if err := cRepo.db.Where(&model.Category{ID: id}).First(&category).Error; err != nil {
		return model.Category{}, err
	}

	return category, nil
}

func (cRepo *categoryRepo) UpdateCategoryByID(id uuid.UUID, category model.Category) (model.Category, error) {
	res := cRepo.db.Where(&model.Category{ID: id}).Clauses(clause.Returning{}).Updates(&category)

	if res.Error != nil {
		return model.Category{}, res.Error
	}

	if res.RowsAffected == 0 {
		return model.Category{}, domain.CategoryErrNotFound
	}

	return category, nil
}

func (cRepo *categoryRepo) DeleteCategoryByID(id uuid.UUID) error {
	if err := cRepo.db.Where(&model.Category{ID: id}).Delete(&model.Category{}).Error; err != nil {
		return err
	}

	return nil
}

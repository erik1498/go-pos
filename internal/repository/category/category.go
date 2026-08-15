package category

import (
	"go-pos/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type categoryRepo struct {
	db *gorm.DB
}

func GetRepository(db *gorm.DB) Repository {
	return &categoryRepo{
		db: db,
	}
}

func (cRepo *categoryRepo) GetCategoryList() ([]model.Category, error) {
	var categoryList = []model.Category{}

	if err := cRepo.db.Where(&model.Category{}).Find(&categoryList).Error; err != nil {
		return nil, err
	}

	return categoryList, nil
}

func (cRepo *categoryRepo) CreateCategory(category model.Category) (model.Category, error) {
	if err := cRepo.db.Create(&category).Error; err != nil {
		return model.Category{}, err
	}

	return category, nil
}

func (cRepo *categoryRepo) GetCategoryById(id uuid.UUID) (model.Category, error) {
	var category model.Category

	if err := cRepo.db.Where(&model.Category{PublicID: id}).First(&category).Error; err != nil {
		return model.Category{}, err
	}

	return category, nil
}

func (cRepo *categoryRepo) UpdateCategoryById(category model.Category, id uuid.UUID) (model.Category, error) {
	if err := cRepo.db.Where(&model.Category{PublicID: id}).Clauses(clause.Returning{}).Updates(&category).Error; err != nil {
		return model.Category{}, err
	}
	return category, nil
}

func (cRepo *categoryRepo) DeleteCategoryByID(id uuid.UUID) error {
	if err := cRepo.db.Where(&model.Category{PublicID: id}).Delete(&model.Category{}).Error; err != nil {
		return err
	}

	return nil
}

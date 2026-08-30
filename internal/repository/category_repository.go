package repository

import (
	"context"
	"errors"
	"fmt"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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

func (cRepo *categoryRepo) GetAll(ctx context.Context, opts domain.QueryOptions) ([]model.Category, int64, error) {
	var categoryList = []model.Category{}
	var totalItems int64

	dbQuery := cRepo.db.WithContext(ctx).Model(&model.Category{})

	if opts.Search != "" {
		dbQuery = dbQuery.Where("name ILIKE ?", "%"+opts.Search+"%")
	}

	for _, f := range opts.Filters {
		queryStr := f.Field + " " + f.Operator + " ?"
		dbQuery = dbQuery.Where(queryStr, f.Value)
	}

	if err := dbQuery.Count(&totalItems).Error; err != nil {
		return nil, 0, fmt.Errorf("[repository][category_repository][GetAll] db query failed: %w", err)
	}

	if err := dbQuery.Order(opts.Sort).Scopes(utils.PaginationScope(opts.Page, opts.Limit)).Find(&categoryList).Error; err != nil {
		return nil, 0, fmt.Errorf("[repository][category_repository][GetAll] db query failed: %w", err)
	}

	return categoryList, totalItems, nil
}

func (cRepo *categoryRepo) Create(ctx context.Context, category model.Category) (model.Category, error) {
	if err := cRepo.db.WithContext(ctx).Create(&category).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "categories_idempotency_key_key" {
				return model.Category{}, fmt.Errorf("[repository][category_repository][Create] err idempotency key: %w", domain.ErrIdempotencyKeyDuplicate)
			}
			return model.Category{}, fmt.Errorf("[repository][category_repository][Create] err duplicate key: %w", err)
		}
		return model.Category{}, fmt.Errorf("[repository][category_repository][Create] db query failed: %w", err)
	}

	return category, nil
}

func (cRepo *categoryRepo) GetByID(ctx context.Context, id uuid.UUID) (model.Category, error) {
	var category model.Category

	if err := cRepo.db.WithContext(ctx).First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Category{}, fmt.Errorf("[repository][category_repository][GetByID] record not found: %w", domain.ErrCategoryNotFound)
		}
		return model.Category{}, fmt.Errorf("[repository][category_repository][GetByID] db query failed: %w", err)
	}

	return category, nil
}

func (cRepo *categoryRepo) UpdateCategoryByID(ctx context.Context, id uuid.UUID, category model.Category) (model.Category, error) {
	res := cRepo.db.WithContext(ctx).Where(&model.Category{ID: id}).Clauses(clause.Returning{}).Updates(&category)

	if res.Error != nil {
		return model.Category{}, fmt.Errorf("[repository][category_repository][UpdateCategoryByID] db query failed: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return model.Category{}, fmt.Errorf("[repository][category_repository][UpdateCategoryByID] record not found: %w", domain.ErrCategoryNotFound)
	}

	return category, nil
}

func (cRepo *categoryRepo) DeleteCategoryByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	err := cRepo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Category{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return fmt.Errorf("[repository][category_repository][DeleteCategoryByID] db query failed: %w", err)
		}

		res := tx.Delete(&model.Category{}, id)

		if res.Error != nil {
			return fmt.Errorf("[repository][category_repository][DeleteCategoryByID] db query failed: %w", res.Error)
		}

		if res.RowsAffected == 0 {
			return fmt.Errorf("[repository][category_repository][DeleteCategoryByID] record not found: %w", domain.ErrCategoryNotFound)
		}

		return nil
	})

	return err
}

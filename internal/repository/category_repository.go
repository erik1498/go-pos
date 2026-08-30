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

type CategoryDAO struct {
	ID             uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid();uniqueIndex;not null;"`
	Name           string         `gorm:"type:varchar(100);not null"`
	IdempotencyKey string         `gorm:"type:varchar(100);uniqueIndex;not null"`
	CreatedBy      uuid.UUID      `gorm:"type:uuid;not null"`
	UpdatedBy      uuid.UUID      `gorm:"type:uuid;not null"`
	DeletedBy      *uuid.UUID     `gorm:"type:uuid"`
	CreatedAt      time.Time      `gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (CategoryDAO) TableName() string {
	return "categories"
}

func (dao *CategoryDAO) ToDomain() domain.Category {
	var deletedAt *time.Time
	if dao.DeletedAt.Valid {
		deletedAt = &dao.DeletedAt.Time
	}
	return domain.Category{
		ID:             dao.ID,
		Name:           dao.Name,
		IdempotencyKey: dao.IdempotencyKey,
		CreatedBy:      dao.CreatedBy,
		UpdatedBy:      dao.UpdatedBy,
		DeletedBy:      dao.DeletedBy,
		CreatedAt:      dao.CreatedAt,
		UpdatedAt:      dao.UpdatedAt,
		DeletedAt:      deletedAt,
	}
}

func FromDomainCategory(c domain.Category) CategoryDAO {
	dao := CategoryDAO{
		ID:             c.ID,
		Name:           c.Name,
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

type categoryRepo struct {
	db *gorm.DB
}

func GetCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepo{
		db: db,
	}
}

func (cRepo *categoryRepo) GetAll(ctx context.Context, opts domain.QueryOptions) ([]domain.Category, int64, error) {
	var categoryList = []domain.Category{}
	var totalItems int64

	dbQuery := cRepo.db.WithContext(ctx).Model(&domain.Category{})

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

func (cRepo *categoryRepo) Create(ctx context.Context, category domain.Category) (domain.Category, error) {
	dao := FromDomainCategory(category)

	if err := cRepo.db.WithContext(ctx).Create(&dao).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "idx_categories_idempotency_key" {
				return domain.Category{}, fmt.Errorf("[repository][category_repository][Create] err idempotency key: %w", domain.ErrIdempotencyKeyDuplicate)
			}
			return domain.Category{}, fmt.Errorf("[repository][category_repository][Create] err duplicate key: %w", err)
		}
		return domain.Category{}, fmt.Errorf("[repository][category_repository][Create] db query failed: %w", err)
	}

	return dao.ToDomain(), nil
}

func (cRepo *categoryRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Category, error) {
	var dao CategoryDAO

	if err := cRepo.db.WithContext(ctx).First(&dao, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Category{}, fmt.Errorf("[repository][category_repository][GetByID] record not found: %w", domain.ErrCategoryNotFound)
		}
		return domain.Category{}, fmt.Errorf("[repository][category_repository][GetByID] db query failed: %w", err)
	}

	return dao.ToDomain(), nil
}

func (cRepo *categoryRepo) UpdateCategoryByID(ctx context.Context, id uuid.UUID, category domain.Category) (domain.Category, error) {
	res := cRepo.db.WithContext(ctx).Where(&domain.Category{ID: id}).Clauses(clause.Returning{}).Updates(&category)

	if res.Error != nil {
		return domain.Category{}, fmt.Errorf("[repository][category_repository][UpdateCategoryByID] db query failed: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return domain.Category{}, fmt.Errorf("[repository][category_repository][UpdateCategoryByID] record not found: %w", domain.ErrCategoryNotFound)
	}

	return category, nil
}

func (cRepo *categoryRepo) DeleteCategoryByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	err := cRepo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.Category{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return fmt.Errorf("[repository][category_repository][DeleteCategoryByID] db query failed: %w", err)
		}

		res := tx.Delete(&domain.Category{}, id)

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

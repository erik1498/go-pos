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
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductDAO struct {
	ID             uuid.UUID       `gorm:"primaryKey;type:uuid;default:gen_random_uuid();uniqueIndex;not null"`
	CategoryID     uuid.UUID       `gorm:"not null"`
	Category       *CategoryDAO    `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Name           string          `gorm:"type:varchar(150);not null;index"`
	SKU            string          `gorm:"type:varchar(50);uniqueIndex:idx_active_sku,where:deleted_at IS NULL;not null"`
	Price          decimal.Decimal `gorm:"type:numeric(12,3);not null;"`
	Stock          decimal.Decimal `gorm:"type:numeric(12,3);default:0;not null;"`
	Taxes          []TaxDAO        `gorm:"many2many:product_taxes;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IdempotencyKey string          `gorm:"type:varchar(100);uniqueIndex;not null"`
	CreatedBy      uuid.UUID       `gorm:"type:uuid;not null"`
	UpdatedBy      uuid.UUID       `gorm:"type:uuid;not null"`
	CreatedAt      time.Time       `gorm:"autoCreateTime"`
	UpdatedAt      time.Time       `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt  `gorm:"index"`
}

func (ProductDAO) TableName() string {
	return "products"
}

func (dao *ProductDAO) ToDomain() domain.Product {
	var deletedAt *time.Time
	if dao.DeletedAt.Valid {
		deletedAt = &dao.DeletedAt.Time
	}

	var categoryPtr *domain.Category
	if dao.Category != nil {
		catDomain := dao.Category.ToDomain()
		categoryPtr = &catDomain
	}

	var taxesList []domain.Tax
	for _, taxDAO := range dao.Taxes {
		taxesList = append(taxesList, taxDAO.ToDomain())
	}

	return domain.Product{
		ID:             dao.ID,
		CategoryID:     dao.CategoryID,
		Category:       categoryPtr,
		Name:           dao.Name,
		SKU:            dao.SKU,
		Price:          dao.Price,
		Stock:          dao.Stock,
		Taxes:          taxesList,
		IdempotencyKey: dao.IdempotencyKey,
		CreatedBy:      dao.CreatedBy,
		UpdatedBy:      dao.UpdatedBy,
		CreatedAt:      dao.CreatedAt,
		UpdatedAt:      dao.UpdatedAt,
		DeletedAt:      deletedAt,
	}
}

func FromDomainProduct(p domain.Product) ProductDAO {
	deletedAt := gorm.DeletedAt{}
	if p.Category.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{Time: *p.Category.DeletedAt, Valid: true}
	}

	category := CategoryDAO{
		ID:             p.Category.ID,
		Name:           p.Category.Name,
		IdempotencyKey: p.Category.IdempotencyKey,
		CreatedBy:      p.Category.CreatedBy,
		UpdatedBy:      p.Category.UpdatedBy,
		DeletedBy:      p.Category.DeletedBy,
		CreatedAt:      p.Category.CreatedAt,
		UpdatedAt:      p.Category.UpdatedAt,
		DeletedAt:      deletedAt,
	}

	var taxes []TaxDAO
	for _, t := range p.Taxes {
		deletedAt = gorm.DeletedAt{}
		if t.DeletedAt != nil {
			deletedAt = gorm.DeletedAt{Time: *p.Category.DeletedAt, Valid: true}
		}
		taxes = append(taxes, TaxDAO{
			ID:             t.ID,
			Name:           t.Name,
			Rate:           t.Rate,
			IsActive:       t.IsActive,
			IdempotencyKey: t.IdempotencyKey,
			CreatedBy:      t.CreatedBy,
			UpdatedBy:      t.UpdatedBy,
			DeletedBy:      t.DeletedBy,
			CreatedAt:      t.CreatedAt,
			UpdatedAt:      t.UpdatedAt,
			DeletedAt:      deletedAt,
		})
	}

	dao := ProductDAO{
		ID:             p.ID,
		CategoryID:     p.CategoryID,
		Category:       &category,
		Name:           p.Name,
		SKU:            p.SKU,
		Price:          p.Price,
		Stock:          p.Stock,
		Taxes:          []TaxDAO{},
		IdempotencyKey: p.IdempotencyKey,
		CreatedBy:      p.CreatedBy,
		UpdatedBy:      p.UpdatedBy,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}

	return dao
}

type productRepository struct {
	db *gorm.DB
}

func GetProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (pRepo *productRepository) GetAll(ctx context.Context, opts domain.QueryOptions) ([]domain.Product, int64, error) {
	var productList []domain.Product
	var totalItems int64

	dbQuery := pRepo.db.Model(&domain.Product{})

	if opts.Search != "" {
		dbQuery = dbQuery.Where("name ILIKE ?", "%"+opts.Search+"%")
	}

	for _, f := range opts.Filters {
		queryStr := f.Field + " " + f.Operator + " ?"
		dbQuery = dbQuery.Where(queryStr, f.Value)
	}

	if err := dbQuery.Count(&totalItems).Error; err != nil {
		return nil, 0, fmt.Errorf("[repository][product_repository][GetAll] db query failed: %w", err)
	}

	if err := dbQuery.Order(opts.Sort).Scopes(utils.PaginationScope(opts.Page, opts.Limit)).Find(&productList).Error; err != nil {
		return nil, 0, fmt.Errorf("[repository][product_repository][GetAll] db query failed: %w", err)
	}

	return productList, totalItems, nil
}

func (pRepo *productRepository) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	if err := pRepo.db.Create(&product).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "idx_products_sku" {
				return domain.Product{}, fmt.Errorf("[repository][product_repository][Create] err sku key: %w", domain.ErrProductSKUIsAlreadyRegistered)
			}
			if pgErr.ConstraintName == "idx_c_idempotency_key" {
				return domain.Product{}, fmt.Errorf("[repository][product_repository][Create] err idempotency key: %w", domain.ErrIdempotencyKeyDuplicate)
			}
			return domain.Product{}, fmt.Errorf("[repository][product_repository][Create] err duplicate key: %w", err)
		}
		return domain.Product{}, fmt.Errorf("[repository][product_repository][Create] db query failed: %w", err)
	}
	return product, nil
}

func (pRepo *productRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	var dao ProductDAO

	if err := pRepo.db.WithContext(ctx).Preload("Taxes").First(&dao, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Product{}, fmt.Errorf("[repository][product_repository][GetByID] record not found: %w", domain.ErrProductNotFound)
		}
		return domain.Product{}, fmt.Errorf("[repository][product_repository][GetByID] db query failed: %w", err)
	}

	return dao.ToDomain(), nil
}

func (pRepo *productRepository) UpdateByID(ctx context.Context, id uuid.UUID, product domain.Product) (domain.Product, error) {
	dao := FromDomainProduct(product)

	err := pRepo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&dao).Clauses(clause.Returning{}).Updates(&dao)
		if res.Error != nil {
			var pgErr *pgconn.PgError
			if errors.As(res.Error, &pgErr) && pgErr.Code == "23505" {
				if pgErr.ConstraintName == "idx_products_sku" {
					return fmt.Errorf("[repository][product_repository][UpdateByID] err sku key: %w", domain.ErrProductSKUIsAlreadyRegistered)
				}
				return fmt.Errorf("[repository][product_repository][UpdateByID] err duplicate key: %w", res.Error)
			}
			return fmt.Errorf("[repository][product_repository][UpdateByID] db query failed: %w", res.Error)
		}

		if res.RowsAffected == 0 {
			return fmt.Errorf("[repository][product_repository][UpdateByID] record not found: %w", domain.ErrProductNotFound)
		}

		if err := tx.Model(&dao).Association("Taxes").Replace(dao.Taxes); err != nil {
			return fmt.Errorf("[repository][product_repository][UpdateByID] db query failed: %w", res.Error)
		}

		return nil
	})

	if err != nil {
		return domain.Product{}, err
	}
	return dao.ToDomain(), nil
}

func (pRepo *productRepository) DeleteByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	err := pRepo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&ProductDAO{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return fmt.Errorf("[repository][product_repository][DeleteByID] db query failed: %w", err)
		}

		res := tx.Delete(&ProductDAO{}, id)

		if res.Error != nil {
			return fmt.Errorf("[repository][product_repository][DeleteByID] db query failed: %w", res.Error)
		}

		if res.RowsAffected == 0 {
			return fmt.Errorf("[repository][product_repository][DeleteByID] record not found: %w", domain.ErrProductNotFound)
		}

		return nil
	})

	return err
}

package repository

import (
	"context"
	"errors"
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
	Stock          decimal.Decimal `gorm:"type:numeric(12,3);default:0;no null;"`
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
	return domain.Product{
		ID:             dao.ID,
		CategoryID:     dao.CategoryID,
		Category:       dao.ToDomain().Category,
		Name:           dao.Name,
		SKU:            dao.SKU,
		Price:          dao.Price,
		Stock:          dao.Stock,
		Taxes:          dao.ToDomain().Taxes,
		IdempotencyKey: dao.IdempotencyKey,
		CreatedBy:      dao.CreatedBy,
		UpdatedBy:      dao.UpdatedBy,
		CreatedAt:      dao.CreatedAt,
		UpdatedAt:      dao.UpdatedAt,
		DeletedAt:      deletedAt,
	}
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
		return nil, 0, err
	}

	err := dbQuery.Order(opts.Sort).Scopes(utils.PaginationScope(opts.Page, opts.Limit)).Find(&productList).Error

	return productList, totalItems, err
}

func (pRepo *productRepository) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	if err := pRepo.db.Create(&product).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.Product{}, domain.ErrProductSKUIsAlreadyRegistered
		}
		return domain.Product{}, err
	}
	return product, nil
}

func (pRepo *productRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	var product domain.Product

	if err := pRepo.db.Preload("Taxes").First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Product{}, domain.ErrProductNotFound
		}
		return domain.Product{}, err
	}

	return product, nil
}

func (pRepo *productRepository) UpdateByID(ctx context.Context, id uuid.UUID, product domain.Product) (domain.Product, error) {
	product.ID = id
	err := pRepo.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&product).Clauses(clause.Returning{}).Updates(&product)
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return domain.ErrProductNotFound
		}

		if err := tx.Model(&product).Association("Taxes").Replace(product.Taxes); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.Product{}, domain.ErrProductSKUIsAlreadyRegistered
		}
		return domain.Product{}, err
	}
	return product, nil
}

func (pRepo *productRepository) DeleteByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	err := pRepo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.Product{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}

		res := tx.Delete(&domain.Product{}, id)

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return domain.ErrProductNotFound
		}

		return nil
	})

	return err
}

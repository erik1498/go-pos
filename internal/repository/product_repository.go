package repository

import (
	"errors"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type productRepository struct {
	db *gorm.DB
}

func GetProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (pRepo *productRepository) GetAll(opts domain.QueryOptions) ([]model.Product, int64, error) {
	var productList []model.Product
	var totalItems int64

	dbQuery := pRepo.db.Model(&model.Product{})

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

func (pRepo *productRepository) Create(product model.Product) (model.Product, error) {
	if err := pRepo.db.Create(&product).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.Product{}, domain.ErrProductSKUIsAlreadyRegistered
		}
		return model.Product{}, err
	}
	return product, nil
}

func (pRepo *productRepository) GetByID(id uuid.UUID) (model.Product, error) {
	var product model.Product

	if err := pRepo.db.Preload("Taxes").First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Product{}, domain.ErrProductNotFound
		}
		return model.Product{}, err
	}

	return product, nil
}

func (pRepo *productRepository) UpdateByID(id uuid.UUID, product model.Product) (model.Product, error) {
	res := pRepo.db.Where(&model.Product{ID: id}).Clauses(clause.Returning{}).Updates(&product)
	if res.Error != nil {
		return model.Product{}, res.Error
	}

	if res.RowsAffected == 0 {
		return model.Product{}, domain.ErrProductNotFound
	}

	return product, nil
}

func (pRepo *productRepository) DeleteByID(id uuid.UUID) error {
	res := pRepo.db.Delete(&model.Product{}, id)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

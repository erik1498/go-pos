package repository

import (
	"errors"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func GetOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (oRepo *orderRepository) GetAll(opts domain.QueryOptions) ([]model.Order, int64, error) {
	var orders []model.Order
	var totalItems int64

	dbQuery := oRepo.db.Model(&model.Order{})

	if opts.Search != "" {
		dbQuery.Where("order_no ILIKE ", "%"+opts.Search+"%")
	}

	for _, f := range opts.Filters {
		queryStr := f.Field + " " + f.Operator + " ?"
		dbQuery.Where(queryStr, f.Value)
	}

	if err := dbQuery.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	if err := dbQuery.Order(opts.Sort).Scopes(utils.PaginationScope(opts.Page, opts.Limit)).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, totalItems, nil

}

func (oRepo *orderRepository) Create(order model.Order) (model.Order, error) {
	err := oRepo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		for _, item := range order.Items {
			res := tx.Model(&model.Product{}).Where("id = ? AND stock >= ?", item.ProductID, item.Qty).UpdateColumn("stock", gorm.Expr("stock - ?", item.Qty))

			if res.Error != nil {
				return res.Error
			}

			if res.RowsAffected == 0 {
				return domain.ErrProductStockIsNotEnough
			}
		}

		return nil
	})
	return order, err
}

func (oRepo *orderRepository) GetByID(id uuid.UUID) (model.Order, error) {
	var order model.Order
	if err := oRepo.db.First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Order{}, domain.ErrOrderNotFound
		}
		return model.Order{}, err
	}
	return order, nil
}

func (oRepo *orderRepository) UpdateByID(id uuid.UUID, order model.Order) (model.Order, error) {
	res := oRepo.db.Where(&model.Order{ID: id}).Updates(&order)

	if res.Error != nil {
		return model.Order{}, res.Error
	}

	if res.RowsAffected == 0 {
		return model.Order{}, domain.ErrOrderNotFound
	}

	return order, nil
}

func (oRepo *orderRepository) DeleteByID(id uuid.UUID) error {
	res := oRepo.db.Delete(&model.Order{}, id)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return domain.ErrOrderNotFound
	}

	return nil
}

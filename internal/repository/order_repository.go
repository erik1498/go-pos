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
)

type PaymentStatus string

type PaymentMethod string

type OrderDAO struct {
	ID                    uuid.UUID       `gorm:"primaryKey;type:uuid;default:gen_random_uuid();uniqueIndex;not null"`
	OrderNo               string          `gorm:"type:varchar(50);uniqueIndex;not null;"`
	MemberID              *uuid.UUID      `gorm:"type:uuid;index"`
	Member                *MemberDAO      `gorm:"foreignKey:MemberID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	PaymentMethod         PaymentMethod   `gorm:"type:varchar(30);not null;default:'CASH'"`
	PaymentStatus         PaymentStatus   `gorm:"type:varchar(30);not null;default:'PENDING'"`
	MidtransTransactionID *string         `gorm:"type:varchar(255);uniqueIndex"`
	SnapToken             *string         `gorm:"type:varchar(255)"`
	PaidAt                *time.Time      `gorm:"type:timestamptz"`
	TotalQty              decimal.Decimal `gorm:"type:numeric(12,3);not null"`
	SubTotal              decimal.Decimal `gorm:"type:numeric(12,3);not null"`
	TotalTax              decimal.Decimal `gorm:"type:numeric(12,3);not null"`
	GrandTotal            decimal.Decimal `gorm:"type:numeric(12,3);not null"`
	IdempotencyKey        string          `gorm:"type:varchar(100);uniqueIndex;not null"`
	CreatedBy             uuid.UUID       `gorm:"type:uuid;not null"`
	CreatedAt             time.Time       `gorm:"autoCreateTime"`
	Items                 []OrderItemDAO  `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}

type OrderItemDAO struct {
	ID             uuid.UUID         `gorm:"primaryKey;type:uuid;default:gen_random_uuid();uniqueIndex;not null"`
	OrderID        uuid.UUID         `gorm:"type:uuid;not null;index"`
	ProductID      uuid.UUID         `gorm:"type:uuid;not null;index"`
	ProductName    string            `gorm:"type:varchar(150);not null"`
	BasePrice      decimal.Decimal   `gorm:"type:numeric(12,3);not null"`
	Qty            decimal.Decimal   `gorm:"type:numeric(12,3);not null"`
	SubTotal       decimal.Decimal   `gorm:"type:numeric(12,3);not null"`
	IdempotencyKey string            `gorm:"type:varchar(100);uniqueIndex;not null"`
	CreatedBy      uuid.UUID         `gorm:"type:uuid;not null"`
	CreatedAt      time.Time         `gorm:"autoCreateTime"`
	AppliedTaxes   []OrderItemTaxDAO `gorm:"foreignKey:OrderItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type OrderItemTaxDAO struct {
	ID             uuid.UUID       `gorm:"primaryKey;type:uuid;default:gen_random_uuid();uniqueIndex;not null"`
	OrderItemID    uuid.UUID       `gorm:"type:uuid;not null;index"`
	TaxID          *uuid.UUID      `gorm:"type:uuid"`
	TaxName        string          `gorm:"type:varchar(50);not null"`
	TaxRate        decimal.Decimal `gorm:"type:numeric(5,2);not null"`
	IdempotencyKey string          `gorm:"type:varchar(100);uniqueIndex;not null"`
	CreatedAt      time.Time       `gorm:"autoCreateTime"`
	CreatedBy      uuid.UUID       `gorm:"type:uuid;not null"`
	TaxAmount      decimal.Decimal `gorm:"type:numeric(12,3);not null"`
}

type orderRepository struct {
	db *gorm.DB
}

func FromDomainOrder(o domain.Order) OrderDAO {
	dao := OrderDAO{
		ID:                    o.ID,
		OrderNo:               o.OrderNo,
		MemberID:              o.MemberID,
		PaymentMethod:         PaymentMethod(o.PaymentMethod),
		PaymentStatus:         PaymentStatus(o.PaymentStatus),
		MidtransTransactionID: o.MidtransTransactionID,
		SnapToken:             o.SnapToken,
		PaidAt:                o.PaidAt,
		TotalQty:              o.TotalQty,
		SubTotal:              o.SubTotal,
		TotalTax:              o.TotalTax,
		GrandTotal:            o.GrandTotal,
		IdempotencyKey:        o.IdempotencyKey,
		CreatedBy:             o.CreatedBy,
		CreatedAt:             o.CreatedAt,
	}
	return dao
}

func (dao *OrderDAO) ToDomain() domain.Order {
	return domain.Order{
		ID:                    dao.CreatedBy,
		OrderNo:               dao.OrderNo,
		MemberID:              dao.MemberID,
		PaymentMethod:         domain.PaymentMethod(dao.PaymentMethod),
		PaymentStatus:         domain.PaymentStatus(dao.PaymentStatus),
		MidtransTransactionID: dao.MidtransTransactionID,
		SnapToken:             dao.SnapToken,
		PaidAt:                dao.PaidAt,
		TotalQty:              dao.TotalQty,
		SubTotal:              dao.SubTotal,
		TotalTax:              dao.TotalTax,
		GrandTotal:            dao.GrandTotal,
		CreatedBy:             dao.CreatedBy,
		CreatedAt:             dao.CreatedAt,
		IdempotencyKey:        dao.IdempotencyKey,
	}
}

func GetOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (oRepo *orderRepository) GetAll(ctx context.Context, opts domain.QueryOptions) ([]domain.Order, int64, error) {
	var daoList []OrderDAO
	var totalItems int64

	dbQuery := oRepo.db.WithContext(ctx).Model(&OrderDAO{})

	if opts.Search != "" {
		dbQuery.Where("order_no ILIKE ", "%"+opts.Search+"%")
	}

	for _, f := range opts.Filters {
		queryStr := f.Field + " " + f.Operator + " ?"
		dbQuery.Where(queryStr, f.Value)
	}

	if err := dbQuery.Count(&totalItems).Error; err != nil {
		return nil, 0, fmt.Errorf("[repository][order_repository][GetAll] db query failed: %w", err)
	}

	if err := dbQuery.Order(opts.Sort).Scopes(utils.PaginationScope(opts.Page, opts.Limit)).Find(&daoList).Error; err != nil {
		return nil, 0, fmt.Errorf("[repository][order_repository][GetAll] db query failed: %w", err)
	}

	var orders []domain.Order
	for _, o := range daoList {
		orders = append(orders, o.ToDomain())
	}

	return orders, totalItems, nil

}

func (oRepo *orderRepository) Create(ctx context.Context, order domain.Order) (domain.Order, error) {
	dao := FromDomainOrder(order)

	err := oRepo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dao).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				if pgErr.ConstraintName == "idx_orders_idempotency_key" {
					return fmt.Errorf("[repository][order_repository][Create] err idempotency key: %w", domain.ErrIdempotencyKeyDuplicate)
				}
				if pgErr.ConstraintName == "idx_orders_order_no" {
					return fmt.Errorf("[repository][order_repository][Create] err order no duplicate: %w", domain.ErrOrderNoIsAlreadyRegistered)
				}
				return fmt.Errorf("[repository][order_repository][Create] err duplicate key: %w", err)
			}
			return fmt.Errorf("[repository][order_repository][Create] db query failed: %w", err)
		}

		for _, item := range order.Items {
			res := tx.Model(&ProductDAO{}).Where("id = ? AND stock >= ?", item.ProductID, item.Qty).UpdateColumn("stock", gorm.Expr("stock - ?", item.Qty))

			if res.Error != nil {
				return fmt.Errorf("[repository][order_repository][Create] db query failed: %w", res.Error)
			}

			if res.RowsAffected == 0 {
				return fmt.Errorf("[repository][order_repository][Create] db query failed: %w", domain.ErrProductStockIsNotEnough)
			}
		}

		return nil
	})

	if err != nil {
		return domain.Order{}, err
	}

	return dao.ToDomain(), nil
}

func (oRepo *orderRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Order, error) {
	var dao OrderDAO
	if err := oRepo.db.WithContext(ctx).First(&dao, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Order{}, fmt.Errorf("[repository][order_repository][GetByID] record not found: %w", domain.ErrOrderNotFound)
		}
		return domain.Order{}, fmt.Errorf("[repository][order_repository][GetByID] db query failed: %w", err)
	}
	return dao.ToDomain(), nil
}

package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "PENDING"
	PaymentStatusPaid     PaymentStatus = "PAID"
	PaymentStatusExpired  PaymentStatus = "EXPIRED"
	PaymentStatusFailed   PaymentStatus = "FAILED"
	PaymentStatusCanceled PaymentStatus = "CANCEL"
)

type PaymentMethod string

const (
	PaymentMethodCash       PaymentMethod = "CASH"
	PaymentMethodMidtransQR PaymentMethod = "MIDTRANS_QR"
	PaymentMethodMidtransVA PaymentMethod = "MIDTRANS_VA"
	PaymentMethodMidtransCC PaymentMethod = "MIDTRANS_CC"
)

type Order struct {
	ID                    uuid.UUID       `gorm:"primaryKey;type:uuid;default:gen_random_uuid();uniqueIndex;not null" json:"id"`
	OrderNo               string          `gorm:"type:varchar(50);uniqueIndex;not null;" json:"order_no"`
	MemberID              *uuid.UUID      `gorm:"type:uuid;index" json:"member_id,omitempty"`
	Member                *Member         `gorm:"foreignKey:MemberID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"member,omitempty"`
	PaymentMethod         PaymentMethod   `gorm:"type:varchar(30);not null;default:'CASH'" json:"payment_method"`
	PaymentStatus         PaymentStatus   `gorm:"type:varchar(30);not null;default:'PENDING'" json:"payment_status"`
	MidtransTransactionID *string         `gorm:"type:varchar(255);uniqueIndex" json:"midtrans_transaction_id,omitempty"`
	SnapToken             *string         `gorm:"type:varchar(255)" json:"snap_token,omitempty"`
	PaidAt                *time.Time      `gorm:"type:timestamptz" json:"paid_at,omitempty"`
	TotalQty              decimal.Decimal `gorm:"type:numeric(12,3);not null" json:"total_qty"`
	SubTotal              decimal.Decimal `gorm:"type:numeric(12,3);not null" json:"sub_total"`
	TotalTax              decimal.Decimal `gorm:"type:numeric(12,3);not null" json:"total_tax"`
	GrandTotal            decimal.Decimal `gorm:"type:numeric(12,3);not null" json:"grand_total"`
	CreatedAt             time.Time       `gorm:"autoCreateTime" json:"created_at"`
	Items                 []OrderItem     `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE" json:"items"`
}

type OrderItem struct {
	ID           uuid.UUID       `gorm:"primaryKey;type:uuid;default:gen_random_uuid();uniqueIndex;not null" json:"id"`
	OrderID      uuid.UUID       `gorm:"type:uuid;not null;index" json:"order_id"`
	ProductID    uuid.UUID       `gorm:"type:uuid;not null;index" json:"product_id"`
	ProductName  string          `gorm:"type:varchar(150);not null" json:"product_name"`
	BasePrice    decimal.Decimal `gorm:"type:numeric(12,3);not null" json:"base_price"`
	Qty          decimal.Decimal `gorm:"type:numeric(12,3);not null" json:"qty"`
	SubTotal     decimal.Decimal `gorm:"type:numeric(12,3);not null" json:"sub_total"`
	AppliedTaxes []OrderItemTax  `gorm:"foreignKey:OrderItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"applied_taxes"`
}

type OrderItemTax struct {
	ID          uuid.UUID       `gorm:"primaryKey;type:uuid;default:gen_random_uuid();uniqueIndex;not null" json:"id"`
	OrderItemID uuid.UUID       `gorm:"type:uuid;not null;index" json:"order_item_id"`
	TaxID       *uuid.UUID      `gorm:"type:uuid" json:"tax_id"`
	TaxName     string          `gorm:"type:varchar(50);not null" json:"tax_name"`
	TaxRate     decimal.Decimal `gorm:"type:numeric(5,2);not null" json:"tax_rate"`
	TaxAmount   decimal.Decimal `gorm:"type:numeric(12,3);not null" json:"tax_amount"`
}

type CreateOrderRequest struct {
	MemberID      *uuid.UUID         `json:"member_id"`
	PaymentMethod PaymentMethod      `json:"payment_method" validate:"required"`
	Items         []OrderItemRequest `json:"items" validate:"required,min=1"`
}

type OrderItemRequest struct {
	ProductID uuid.UUID       `json:"product_id" validate:"required"`
	Qty       decimal.Decimal `json:"qty" validate:"required,gt=0"`
}

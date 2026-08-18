package model

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "PENDING"
	PaymentStatusPaid     PaymentStatus = "PAID"
	PaymentStatusExpired  PaymentStatus = "EXPIRED"
	PaymentStatusFailed   PaymentStatus = "FAILED"
	PaymentStatusCanceled PaymentStatus = "CANCEL"
)

type Order struct {
	ID                    uuid.UUID     `gorm:"primaryKey;type:uuid;default:gen_random_uuid();uniqueIndex;not null" json:"id"`
	OrderNo               string        `gorm:"type:varchar(50);uniqueIndex;not null;" json:"order_no"`
	MemberID              uuid.UUID     `gorm:"type:uuid;index" json:"member_id"`
	Member                *Member       `gorm:"foreignKey:MemberID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"member,omitempty"`
	PaymentMethod         string        `gorm:"type:varchar(30);not null;default:'CASH'" json:"payment_method"`
	PaymentStatus         PaymentStatus `gorm:"type:varchar(30);not null;default:'PENDING'" json:"payment_status"`
	MidtransTransactionID *string       `gorm:"type:varchar(255);uniqueIndex" json:"midtrans_transaction_id,omitempty"`
	SnapToken             *string       `gorm:"type:varchar(255)" json:"snap_token,omitempty"`
	PaidAt                *time.Time    `gorm:"type:timestamptz" json:"paid_at,omitempty"`
	TotalQty              float64       `gorm:"type:numeric(12,3);not null" json:"total_qty"`
	SubTotal              float64       `gorm:"type:numeric(12,3);not null" json:"sub_total"`
	TotalTax              float64       `gorm:"type:numeric(12,3);not null" json:"total_tax"`
	GrandTotal            float64       `gorm:"type:numeric(12,3);not null" json:"grand_total"`
	CreatedAt             time.Time     `gorm:"autoCreateTime" json:"created_at"`
	Items                 []OrderItem   `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE" json:"items"`
}

type OrderItem struct {
	ID           uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid();uniqueIndex;not null" json:"id"`
	OrderID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"order_id"`
	ProductID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"product_id"`
	ProductName  string         `gorm:"type:varchar(150);not null" json:"product_name"`
	BasePrice    float64        `gorm:"type:numeric(12,3);not null" json:"base_price"`
	Qty          float64        `gorm:"type:numeric(12,3);not null" json:"qty"`
	SubTotal     float64        `gorm:"type:numeric(12,3);not null" json:"sub_total"`
	AppliedTaxes []OrderItemTax `gorm:"foreignKey:OrderItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"applied_taxes"`
}

type OrderItemTax struct {
	ID          uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid();uniqueIndex;not null" json:"id"`
	OrderItemID uuid.UUID  `gorm:"type:uuid;not null;index" json:"order_item_id"`
	TaxID       *uuid.UUID `gorm:"type:uuid" json:"tax_id"`
	TaxName     string     `gorm:"type:varchar(50);not null" json:"tax_name"`
	TaxRate     float64    `gorm:"type:numeric(5,2);not null" json:"tax_rate"`
	TaxAmount   float64    `gorm:"type:numeric(12,3);not null" json:"tax_amount"`
}

type OrderRequest struct {
	OrderNo       string `json:"order_no"`
	MemberID      string `json:"member_id"`
	PaymentMethod string `json:"payment_method"`
}

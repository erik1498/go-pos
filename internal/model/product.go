package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	ID         uuid.UUID       `gorm:"primaryKey;type:uuid;default:gen_random_uuid();uniqueIndex;not null" json:"id"`
	CategoryID uuid.UUID       `gorm:"not null" json:"category_id"`
	Category   *Category       `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"category,omitempty"`
	Name       string          `gorm:"type:varchar(150);not null;index" json:"name"`
	SKU        string          `gorm:"type:varchar(50);uniqueIndex:idx_active_sku,where:deleted_at IS NULL;not null" json:"sku"`
	Price      decimal.Decimal `gorm:"type:numeric(12,3);not null;" json:"price"`
	Stock      decimal.Decimal `gorm:"type:numeric(12,3);default:0;no null;" json:"stock"`
	Taxes      []Tax           `gorm:"many2many:product_taxes;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"taxes,omitempty"`
	CreatedAt  time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"`
}

type ProductRequest struct {
	Name       string          `json:"name"`
	SKU        string          `json:"sku"`
	Price      decimal.Decimal `json:"price"`
	CategoryID string          `json:"category_id"`
	Tax        []Tax           `json:"tax"`
}

type ProductTaxRequest struct {
	TaxID string `json:"tax_id"`
}

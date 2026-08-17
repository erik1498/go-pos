package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID         uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid();uniqueIndex;not null" json:"id"`
	CategoryID uuid.UUID      `gorm:"not null" json:"category_id"`
	Category   *Category      `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"category,omitempty"`
	Name       string         `gorm:"type:varchar(150);not null;index" json:"name"`
	SKU        string         `gorm:"type:varchar(50);uniqueIndex:idx_active_sku,where deleted_at IS NULL;not null" json:"sku"`
	Price      float64        `gorm:"type:numeric(12,3);not null;" json:"price"`
	Stock      float64        `gorm:"type:numeric(12,3);default:0;no null;" json:"stock"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type ProductRequest struct {
	Name       string  `json:"name"`
	SKU        string  `json:"sku"`
	Price      float64 `json:"price"`
	CategoryID string  `json:"category_id"`
}

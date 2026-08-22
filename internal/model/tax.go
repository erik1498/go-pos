package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Tax struct {
	ID        uuid.UUID       `gorm:"primaryKey;type:uuid;default:gen_random_uuid();uniqueIndex;not null" json:"id"`
	Name      string          `gorm:"type:varchar(50);uniqueIndex:idx_active_tax_name,where:deleted_at IS NULL;not null" json:"name"`
	Rate      decimal.Decimal `gorm:"type:numeric(5,2);not null" json:"rate"`
	IsActive  bool            `gorm:"type:boolean;default:true" json:"is_active"`
	CreatedAt time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

type TaxRequest struct {
	Name string          `json:"name"`
	Rate decimal.Decimal `json:"rate"`
}

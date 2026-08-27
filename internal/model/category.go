package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID             uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid();uniqueIndex;not null;" json:"id"`
	Name           string         `gorm:"type:varchar(100);not null" json:"name"`
	IdempotencyKey string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"idempotency_key"`
	CreatedBy      uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	UpdatedBy      uuid.UUID      `gorm:"type:uuid;not null" json:"updated_by"`
	DeletedBy      *uuid.UUID     `gorm:"type:uuid" json:"deleted_by,omitempty"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type CategoryRequest struct {
	Name string `json:"name"`
}

package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        int64          `gorm:"primaryKey;autoIncrement;" json:"-"`
	PublicID  uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;not null;" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CategoryRequest struct {
	Name string `json:"name"`
}

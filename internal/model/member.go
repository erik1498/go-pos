package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Member struct {
	ID             uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid();uniqueIndex;not null" json:"id"`
	MemberCode     string         `gorm:"uniqueIndex;not null;type:varchar(50)" json:"member_code"`
	NameEncrypted  []byte         `gorm:"type:bytea;not null" json:"-"`
	PhoneEncrypted []byte         `gorm:"type:bytea;not null" json:"-"`
	EmailEncrypted []byte         `gorm:"type:bytea;not null" json:"-"`
	PhoneBIndex    string         `gorm:"type:varchar(64);uniqueIndex:idx_active_phone_bindex,where:deleted_at IS NULL;not null" json:"-"`
	EmailBIndex    *string        `gorm:"type:varchar(64);uniqueIndex:idx_active_email_bindex,where:deleted_at IS NULL" json:"-"`
	Name           string         `gorm:"-" json:"name"`
	Phone          string         `gorm:"-" json:"phone"`
	Email          string         `gorm:"-" json:"email,omitempty"`
	Points         int            `gorm:"not null;default:0;" json:"points"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type MemberRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

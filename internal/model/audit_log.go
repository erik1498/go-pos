package model

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ActorID   uuid.UUID `gorm:"type:uuid;index;not null" json:"actor_id"`
	ActorRole string    `gorm:"type:varchar(50);not null" json:"actor_role"`
	Action    string    `gorm:"type:varchar(20);not null" json:"action"`
	Entity    string    `gorm:"type:varchar(50);not null" json:"entity"`
	EntityID  string    `gorm:"type:varchar(255);index;not null" json:"entity_id"`
	OldValues string    `gorm:"type:jsonb" json:"old_values"`
	NewValues string    `gorm:"type:jsonb" json:"new_values"`
	IPAddress string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent string    `gorm:"type:text" json:"user_agent"`
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

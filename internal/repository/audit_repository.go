package repository

import (
	"context"
	"go-pos/internal/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLogDAO struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ActorID   uuid.UUID `gorm:"type:uuid;index;not null"`
	ActorRole string    `gorm:"type:varchar(50);not null"`
	Action    string    `gorm:"type:varchar(20);not null"`
	Entity    string    `gorm:"type:varchar(50);not null" json:"entity"`
	EntityID  string    `gorm:"type:varchar(255);index;not null" json:"entity_id"`
	OldValues string    `gorm:"type:jsonb" json:"old_values"`
	NewValues string    `gorm:"type:jsonb" json:"new_values"`
	IPAddress string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent string    `gorm:"type:text" json:"user_agent"`
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

func (AuditLogDAO) TableName() string {
	return "audit_logs"
}

func FromDomainAuditLog(a domain.AuditLog) AuditLogDAO {
	dao := AuditLogDAO{
		ID:        a.ID,
		ActorID:   a.ActorID,
		ActorRole: a.ActorRole,
		Action:    a.Action,
		Entity:    a.Entity,
		EntityID:  a.EntityID,
		OldValues: a.OldValues,
		NewValues: a.NewValues,
		IPAddress: a.IPAddress,
		UserAgent: a.UserAgent,
		CreatedAt: a.CreatedAt,
	}
	return dao
}

type auditRepository struct {
	db *gorm.DB
}

func GetAuditLogRepository(db *gorm.DB) domain.AuditLogRepository {
	return &auditRepository{
		db: db,
	}
}

func (alRepo *auditRepository) Create(ctx context.Context, logData domain.AuditLog) error {
	dao := FromDomainAuditLog(logData)

	if err := alRepo.db.WithContext(ctx).Create(&dao).Error; err != nil {
		return err
	}

	return nil
}

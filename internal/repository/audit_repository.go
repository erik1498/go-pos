package repository

import (
	"context"
	"go-pos/internal/domain"
	"go-pos/internal/model"

	"gorm.io/gorm"
)

type auditRepository struct {
	db *gorm.DB
}

func GetAuditLogRepository(db *gorm.DB) domain.AuditLogRepository {
	return &auditRepository{
		db: db,
	}
}

func (al *auditRepository) Create(ctx context.Context, logData model.AuditLog) {
	al.db.Create(&logData)
}

package database

import (
	"go-pos/internal/model"
	"go-pos/internal/repository"

	"gorm.io/gorm"
)

func seedDB(db *gorm.DB) {
	db.AutoMigrate(
		&repository.CategoryDAO{},
		&repository.ProductDAO{},
		&repository.TaxDAO{},
		&model.AuditLog{},
	)
}

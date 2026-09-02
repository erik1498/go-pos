package database

import (
	"go-pos/internal/repository"

	"gorm.io/gorm"
)

func seedDB(db *gorm.DB) {
	db.AutoMigrate(
		&repository.MemberDAO{},
		&repository.CategoryDAO{},
		&repository.TaxDAO{},
		&repository.ProductDAO{},
		&repository.AuditLogDAO{},
		&repository.UserDAO{},
	)
}

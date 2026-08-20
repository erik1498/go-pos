package database

import (
	"go-pos/internal/model"

	"gorm.io/gorm"
)

func seedDB(db *gorm.DB) {
	db.AutoMigrate(&model.Category{}, &model.Product{}, &model.Member{}, &model.Order{}, &model.OrderItem{}, &model.OrderItemTax{}, &model.Tax{})
}

package database

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("POSTGRES: FAILED TO CONNECT TO DATABASE, ERR: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("POSTGRES: FAILED TO EXTRACT SQL DB INSTANCE, ERR: %v", err)
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	seedDB(db)

	log.Println("POSTGRES: CONNECTED & POOLING CONFIGURED")
	return db
}

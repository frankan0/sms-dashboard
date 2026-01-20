package database

import (
	"log"
	"sms-dashboard/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Migrate the schema
	err = DB.AutoMigrate(&model.SMS{}, &model.User{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}

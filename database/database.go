package database

import (
	"github.com/simon4545/goshop/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	DB.AutoMigrate(
		&models.ProductCategory{},
		&models.Product{},
		&models.NewsCategory{},
		&models.News{},
		&models.ProductCategory{})
}

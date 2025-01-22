package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductCategory struct {
	gorm.Model
	Name string `json:"name"`
}

type Product struct {
	ID                uint `gorm:"primarykey"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Price             float64 `json:"price"`
	ProductCategoryID uint    `json:"category_id"`
	ProductImages     string  `json:"images" gorm:"column:images"`
	SKUs              string  `json:"skus" gorm:"column:skus"`
}

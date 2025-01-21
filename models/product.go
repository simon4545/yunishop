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
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	Price             float64        `json:"price"`
	ProductCategoryID uint           `json:"category_id"`
	ProductImages     []ProductImage `json:"images"`
	SKUs              []SKU          `json:"skus"`
}

type SKU struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ProductID uint    `json:"product_id"`
	Variant   string  `json:"variant"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
}

type ProductImage struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ProductID uint   `json:"product_id"`
	URL       string `json:"url"`
}

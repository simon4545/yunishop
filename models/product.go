package models

import "gorm.io/gorm"

type ProductCategory struct {
	gorm.Model
	Name string `json:"name"`
}

type Product struct {
	gorm.Model
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       float64         `json:"price"`
	CategoryID  uint            `json:"category_id"`
	Category    ProductCategory `gorm:"foreignKey:CategoryID"`
	Images      []ProductImage  `gorm:"foreignKey:ProductID"`
}

type SKU struct {
	gorm.Model
	ProductID uint    `json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID"`
	Variant   string  `json:"variant"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
}

type ProductImage struct {
	gorm.Model
	ProductID uint   `json:"product_id"`
	URL       string `json:"url"`
}

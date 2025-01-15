package models

import "gorm.io/gorm"

type NewsCategory struct {
	gorm.Model
	Name string `json:"name"`
}

type News struct {
	gorm.Model
	Title      string       `json:"title"`
	Content    string       `json:"content"`
	CategoryID uint         `json:"category_id"`
	Category   NewsCategory `gorm:"foreignKey:CategoryID"`
}

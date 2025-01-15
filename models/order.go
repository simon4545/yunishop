package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	Status     string `json:"status"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Address    string `json:"address"`
	PostalCode string `json:"postal_code"`
}

type Item struct {
	gorm.Model
	Order     Order
	OrderID   uint
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

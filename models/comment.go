package models

type ProductReview struct {
	ID        uint     `gorm:"primarykey"`
	SKU       string   `json:"sku"`
	Comment   string   `json:"comment"`
	Images    []string `json:"images"`
	Rating    int      `json:"rating"`
	UserName  string   `json:"user_name"`
	UserID    string   `json:"user_id"`
	CreatedAt int64    `json:"created_at"`
}

package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 定義來源表結構（對應 a.sqlite）
type ItemView struct {
	Title   string  `gorm:"column:標題"`
	Content string  `gorm:"column:說明"`
	Price   float64 `gorm:"column:直購價"`
	Images  string  `gorm:"column:圖片"`
	Skus    string  `gorm:"column:自訂規格"`
	URL     string  `gorm:"column:網址"`
}

// 定義目標表結構（對應 b.sqlite）
type Product struct {
	Title   string  `gorm:"column:title"`
	Content string  `gorm:"column:content"`
	Price   float64 `gorm:"column:price"`
	Images  string  `gorm:"column:images"`
	Skus    string  `gorm:"column:skus"`
	URL     string  `gorm:"column:url"`
}

func main() {
	// 連接 a.sqlite
	srcDB, err := gorm.Open(sqlite.Open("../ruten.sqlite"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to source database: %v", err)
	}

	// 連接 b.sqlite
	dstDB, err := gorm.Open(sqlite.Open("../test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to destination database: %v", err)
	}

	// 查詢所有來源數據
	var items []ItemView
	if err := srcDB.Table("itemview").Find(&items).Error; err != nil {
		log.Fatalf("Failed to fetch data from itemview: %v", err)
	}

	// 插入到目標表
	for _, item := range items {
		product := Product{
			Title:   item.Title,
			Content: item.Content,
			Price:   item.Price,
			Images:  item.Images,
			Skus:    item.Skus,
			URL:     item.URL,
		}
		if err := dstDB.Table("products").Create(&product).Error; err != nil {
			log.Printf("Failed to insert record: %v", err)
		}
	}

	log.Println("Data migration completed successfully!")
}

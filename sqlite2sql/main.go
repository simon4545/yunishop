package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 定義來源表結構（對應 a.sqlite）
type ItemView struct {
	Title   string `gorm:"column:標題"`
	Content string `gorm:"column:說明"`
	Price   string `gorm:"column:直購價"`
	Images  string `gorm:"column:圖片"`
	Skus    string `gorm:"column:自訂規格"`
	URL     string `gorm:"column:網址"`
}

// 定義目標表結構（對應 b.sqlite）
type Product struct {
	Name    string  `gorm:"column:name"`
	Content string  `gorm:"column:description"`
	Price   float64 `gorm:"column:price"`
	Images  string  `gorm:"column:images"`
	Skus    string  `gorm:"column:skus"`
	URL     string  `gorm:"column:url"`
}

func stringToFloat(s string) float64 {
	if s == "" {
		return 0
	}
	// 将字符串转换为 float64
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Println("Error parsing string:", err)
		return 0
	}
	return value
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
	var count int64
	if err := srcDB.Table("itemview").Count(&count).Error; err != nil {
		log.Fatalf("Failed to count records from itemview: %v", err)
	}

	log.Printf("Total records in a.sqlite: %d", count)
	// 逐行處理數據
	rows, err := srcDB.Table("itemview").Rows()
	if err != nil {
		log.Fatalf("Failed to fetch data from itemview: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item ItemView
		if err := srcDB.ScanRows(rows, &item); err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}
		item.Images = strings.ReplaceAll(item.Images, "\\", "/")
		item.Images = strings.ReplaceAll(item.Images, `C:/lutian/PicBackup`, "")
		product := Product{
			Name:    item.Title,
			Content: item.Content,
			Price:   stringToFloat(item.Price),
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

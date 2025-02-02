package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SKU struct {
	Image string            `json:"image"`
	Name  map[string]string `json:"name"`
}

var temp = `[
    {
        "name": {
            "zh_hk": "規格",
            "zh_cn": "規格"
        },
        "values": [],
		"isImage": false
}`

// 定义结构体
type SQLiteProduct struct {
	ID          uint    `gorm:"primaryKey"`
	Title       string  `gorm:"column:標題"`
	Description string  `gorm:"column:說明"`
	Price       float64 `gorm:"column:直購價"`
	Image       string  `gorm:"column:圖片"`
	CustomSpec  string  `gorm:"column:自訂規格"`
}

func (SQLiteProduct) TableName() string {
	return "itemview"
}

type MySQLProduct struct {
	ID         uint    `gorm:"primaryKey"`
	Price      float64 `gorm:"column:price"`
	Image      string  `gorm:"column:images"`
	CustomSpec string  `gorm:"column:variables"`
}

func (MySQLProduct) TableName() string {
	return "products_1"
}
func main() {
	// 连接 SQLite 数据库
	sqliteDB, err := gorm.Open(sqlite.Open("../ruten.sqlite"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to SQLite: %v", err)
	}

	// 连接 MySQL 数据库
	dsn := "sql_192_168_1_30:3839b7857eb1f@tcp(192.168.1.30:3306)/sql_192_168_1_30?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// 确保 MySQL 表存在
	err = mysqlDB.AutoMigrate(&MySQLProduct{})
	if err != nil {
		log.Fatalf("Failed to migrate MySQL schema: %v", err)
	}
	var count int64
	if err := sqliteDB.Table("itemview").Count(&count).Error; err != nil {
		log.Fatalf("Failed to count records from itemview: %v", err)
	}

	log.Printf("Total records in a.sqlite: %d", count)
	// 逐行處理數據
	rows, err := sqliteDB.Table("itemview").Rows()
	if err != nil {
		log.Fatalf("Failed to fetch data from itemview: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item SQLiteProduct
		if err := sqliteDB.ScanRows(rows, &item); err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}

		// 处理图片字段
		item.Image = strings.ReplaceAll(item.Image, "\\", "/")
		item.Image = strings.ReplaceAll(item.Image, `C:/lutian/PicBackup`, "")
		images := strings.Split(item.Image, "|")
		jsonData, err := json.Marshal(images)
		if err != nil {
			fmt.Println("Error serializing slice:", err)
			return
		}

		// 转换为字符串输出
		jsonString := string(jsonData)

		skus := gjson.Get(item.CustomSpec, "1").Map()
		skulist := []SKU{}
		for k, _ := range skus {
			lang := map[string]string{"zh_hk": "規格", "zh_cn": "规格"}
			lang["zh_cn"] = k
			lang["zh_hk"] = k
			skulist = append(skulist, SKU{Name: lang, Image: ""})
		}
		value, _ := sjson.Set(temp, "values", skulist)
		fmt.Println(value)
		product := MySQLProduct{
			Price:      item.Price,
			Image:      jsonString,
			CustomSpec: value,
		}

		if err := mysqlDB.Create(&product).Error; err != nil {
			log.Printf("Failed to insert product into MySQL: %v", err)
		}
		log.Printf("Inserted product: %s", product.ID)
	}

	fmt.Println("Data migration completed successfully!")
}

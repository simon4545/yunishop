package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/simon4545/goshop/database"
	"github.com/simon4545/goshop/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB1() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("product.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.ProductCategory{}, &models.ProductImage{}, &models.Product{}, &models.SKU{})
	return db
}

func TestAddProduct(t *testing.T) {
	db := setupTestDB1()
	database.DB = db
	defer t.Cleanup(func() {
		// Drop all tables
		db.Migrator().DropTable(&models.ProductCategory{}, &models.ProductImage{}, &models.Product{}, &models.SKU{})
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("product.db")
	})
	e := echo.New()
	reqBody := `{
		"name": "Test News",
		"description": "Test Content",
		"price": 1,
		"category_id": 1,
		"images": [
			{
				"product_id": 1,
				"url": "a.jpg"
			}
		],
		"skus": [
			{
				"product_id": 1,
				"variant": "red|large",
				"price": 1111,
				"stock": 123
			}
		]
	}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, AddProduct(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
func CreateProduct(db *gorm.DB) {
	product := models.Product{
		Name:              "test",
		Description:       "test",
		ProductCategoryID: 1,
		ProductImages: []models.ProductImage{
			models.ProductImage{
				ProductID: 1,
				URL:       "a.jpe",
			},
		},
		SKUs: []models.SKU{
			models.SKU{
				ProductID: 1,
				Variant:   "red|large",
				Price:     1212,
				Stock:     123,
			},
		},
	}

	db.Create(&product)
}
func TestDeleteProduct(t *testing.T) {
	db := setupTestDB1()
	database.DB = db
	defer t.Cleanup(func() {
		// Drop all tables
		db.Migrator().DropTable(&models.ProductCategory{}, &models.ProductImage{}, &models.Product{}, &models.SKU{})
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("product.db")
	})

	CreateProduct(db)

	e := echo.New()

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, DeleteProduct(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
func TestUpdateProduct(t *testing.T) {
	db := setupTestDB1()
	database.DB = db
	// defer t.Cleanup(func() {
	// 	db.Migrator().DropTable(&models.ProductCategory{}, &models.ProductImage{}, &models.Product{}, &models.SKU{})
	// 	sqlDB, _ := db.DB()
	// 	sqlDB.Close()
	// 	os.Remove("product.db")
	// })

	CreateProduct(db)

	e := echo.New()
	reqBody := `{
		"name": "Test News1",
		"description": "Test Content1",
		"price": 2,
		"category_id": 2,
		"images": [
			{
				"id":1,
				"product_id": 1,
				"url": "b.jpg"
			}
		],
		"skus": [
			{
				"id":1,
				"product_id": 1,
				"variant": "blue|large",
				"price": 2222,
				"stock": 123
			}
		]
	}`
	req := httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, UpdateProduct(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

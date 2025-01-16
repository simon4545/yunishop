package handlers

import (
	"net/http"
	"net/http/httptest"
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
	db := setupTestDB()
	database.DB = db

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

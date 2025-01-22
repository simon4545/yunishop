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
	db.AutoMigrate(&models.ProductCategory{}, &models.Product{})
	return db
}

func TestAddProduct(t *testing.T) {
	db := setupTestDB1()
	database.DB = db
	defer t.Cleanup(func() {
		// Drop all tables
		db.Migrator().DropTable(&models.ProductCategory{}, &models.Product{})
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
		"images": "b.jpg",
		"skus": ""
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
		ProductImages:     `b.jpg`,
		SKUs: `[
    {
        "dataRows": [
            {
                "name": "颜色分类",
                "id": "1627207",
                "specValue": [
                    {
                        "id": "32123795920",
                        "name": "VS1053模块（背面排针）",
                        "simage": "https://gw.alicdn.com/bao/uploaded/i3/2142287760/O1CN01sT47ey27C85G0BzaT_!!2142287760.jpg"
                    },
                    {
                        "id": "32123795921",
                        "name": "VS1053模块（正面排针）",
                        "simage": "https://gw.alicdn.com/bao/uploaded/i4/2142287760/O1CN01lIAHYZ27C8QGD42Rb_!!2142287760.jpg"
                    },
                    {
                        "id": "3432944",
                        "name": "小音箱",
                        "simage": "https://gw.alicdn.com/bao/uploaded/i4/2142287760/O1CN01S1Pcx327C85Gnl4d3_!!2142287760.jpg"
                    },
                    {
                        "id": "32123795922",
                        "name": "VS1053模块（背面排针）+小音箱",
                        "simage": "https://gw.alicdn.com/bao/uploaded/i1/2142287760/O1CN01F8ssvL27C855Hbcej_!!2142287760.jpg"
                    },
                    {
                        "id": "32123795923",
                        "name": "VS1053模块（正面排针）+小音箱",
                        "simage": "https://gw.alicdn.com/bao/uploaded/i1/2142287760/O1CN01kCc4LQ27C8QMVVCO6_!!2142287760.jpg"
                    }
                ]
            }
        ]
    },
    {
        "VS1053模块（背面排针）": {
            "originalQty": 200,
            "price": 49,
            "skuId": "4584031382428"
        },
        "VS1053模块（正面排针）": {
            "originalQty": 200,
            "price": 49,
            "skuId": "5560890027360"
        },
        "小音箱": {
            "originalQty": 200,
            "price": 19,
            "skuId": "4584031382427"
        },
        "VS1053模块（背面排针）+小音箱": {
            "originalQty": 200,
            "price": 68,
            "skuId": "4584031382429"
        },
        "VS1053模块（正面排针）+小音箱": {
            "originalQty": 200,
            "price": 68,
            "skuId": "5560890027361"
        }
    }
]`,
	}

	db.Create(&product)
}
func TestDeleteProduct(t *testing.T) {
	db := setupTestDB1()
	database.DB = db
	defer t.Cleanup(func() {
		// Drop all tables
		db.Migrator().DropTable(&models.ProductCategory{}, &models.Product{})
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
	defer t.Cleanup(func() {
		db.Migrator().DropTable(&models.ProductCategory{}, &models.Product{})
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("product.db")
	})

	CreateProduct(db)

	e := echo.New()
	reqBody := `{
		"name": "Test News1",
		"description": "Test Content1",
		"price": 2,
		"category_id": 2,
		"images": "a.jpg|b.jpg",
		"skus": "[{\"dataRows\":[{\"name\":\"颜色分类\",\"id\":\"1627207\",\"specValue\":[{\"id\":\"32123795920\",\"name\":\"VS1053模块（背面排针）\",\"simage\":\"https://gw.alicdn.com/bao/uploaded/i3/2142287760/O1CN01sT47ey27C85G0BzaT_!!2142287760.jpg\"},{\"id\":\"32123795921\",\"name\":\"VS1053模块（正面排针）\",\"simage\":\"https://gw.alicdn.com/bao/uploaded/i4/2142287760/O1CN01lIAHYZ27C8QGD42Rb_!!2142287760.jpg\"},{\"id\":\"3432944\",\"name\":\"小音箱\",\"simage\":\"https://gw.alicdn.com/bao/uploaded/i4/2142287760/O1CN01S1Pcx327C85Gnl4d3_!!2142287760.jpg\"},{\"id\":\"32123795922\",\"name\":\"VS1053模块（背面排针）+小音箱\",\"simage\":\"https://gw.alicdn.com/bao/uploaded/i1/2142287760/O1CN01F8ssvL27C855Hbcej_!!2142287760.jpg\"},{\"id\":\"32123795923\",\"name\":\"VS1053模块（正面排针）+小音箱\",\"simage\":\"https://gw.alicdn.com/bao/uploaded/i1/2142287760/O1CN01kCc4LQ27C8QMVVCO6_!!2142287760.jpg\"}]}]},{\"VS1053模块（背面排针）\":{\"originalQty\":200,\"price\":49,\"skuId\":\"4584031382428\"},\"VS1053模块（正面排针）\":{\"originalQty\":200,\"price\":49,\"skuId\":\"5560890027360\"},\"小音箱\":{\"originalQty\":200,\"price\":19,\"skuId\":\"4584031382427\"},\"VS1053模块（背面排针）+小音箱\":{\"originalQty\":200,\"price\":68,\"skuId\":\"4584031382429\"},\"VS1053模块（正面排针）+小音箱\":{\"originalQty\":200,\"price\":68,\"skuId\":\"5560890027361\"}}]"
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

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

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test111.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.NewsCategory{}, &models.News{})
	return db
}

func TestGetNewsCategories(t *testing.T) {
	db := setupTestDB()
	database.DB = db
	defer db.Migrator().DropTable(&models.News{})
	defer db.Migrator().DropTable(&models.NewsCategory{})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/news/categories", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, GetNewsCategories(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestCreateNewsCategory(t *testing.T) {
	db := setupTestDB()
	database.DB = db
	// defer db.Migrator().DropTable(&models.News{})
	// defer db.Migrator().DropTable(&models.NewsCategory{})

	category := models.NewsCategory{Name: "Test Category"}
	db.Create(&category)

	e := echo.New()
	reqBody := `{"name": "Test Category"}`
	req := httptest.NewRequest(http.MethodPost, "/news/categories", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, CreateNewsCategory(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestDeleteNewsCategory(t *testing.T) {
	db := setupTestDB()
	database.DB = db

	// defer db.Migrator().DropTable(&models.News{})
	// defer db.Migrator().DropTable(&models.NewsCategory{})

	category := models.NewsCategory{Name: "Test Category"}
	db.Create(&category)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/news/categories/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, DeleteNewsCategory(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestUpdateNewsCategory(t *testing.T) {
	db := setupTestDB()
	database.DB = db
	defer db.Migrator().DropTable(&models.News{})
	defer db.Migrator().DropTable(&models.NewsCategory{})

	category := models.NewsCategory{Name: "Test Category"}
	db.Create(&category)

	e := echo.New()
	reqBody := `{"name": "Updated Category"}`
	req := httptest.NewRequest(http.MethodPut, "/news/categories/1", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, UpdateNewsCategory(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestCreateNews(t *testing.T) {
	db := setupTestDB()
	database.DB = db
	defer db.Migrator().DropTable(&models.News{})
	defer db.Migrator().DropTable(&models.NewsCategory{})

	category := models.NewsCategory{Name: "Test Category"}
	db.Create(&category)

	e := echo.New()
	reqBody := `{"title": "Test News", "content": "Test Content","category_id":1}`
	req := httptest.NewRequest(http.MethodPost, "/news", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, CreateNews(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestUpdateNews(t *testing.T) {
	db := setupTestDB()
	database.DB = db
	defer db.Migrator().DropTable(&models.News{})

	news := models.News{Title: "Test News", Content: "Test Content"}
	db.Create(&news)

	e := echo.New()
	reqBody := `{"title": "Updated News", "content": "Updated Content"}`
	req := httptest.NewRequest(http.MethodPut, "/news/1", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, UpdateNews(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestDeleteNews(t *testing.T) {
	db := setupTestDB()
	database.DB = db
	// defer db.Migrator().DropTable(&models.News{})
	category := models.NewsCategory{Name: "Test Category"}
	db.Create(&category)

	news := models.News{Title: "Test", Content: "Content", Category: category}
	db.Create(&news)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/news/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, DeleteNews(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestGetNews(t *testing.T) {
	db := setupTestDB()
	database.DB = db
	defer db.Migrator().DropTable(&models.News{})

	news := models.News{Title: "Test News", Content: "Test Content"}
	db.Create(&news)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/news/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, GetNews(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

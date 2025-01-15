package handlers

import (
	"net/http"
	"strconv"

	"github.com/simon4545/goshop/database"
	"github.com/simon4545/goshop/models"

	"github.com/labstack/echo/v4"
)

func GetNewsCategories(c echo.Context) error {
	var categories []models.NewsCategory
	database.DB.Find(&categories)
	return c.JSON(http.StatusOK, categories)
}
func CreateNewsCategory(c echo.Context) error {
	var category models.NewsCategory
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	database.DB.Create(&category)
	return c.JSON(http.StatusCreated, category)
}

func DeleteNewsCategory(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var category models.NewsCategory
	database.DB.First(&category, id)
	if category.ID == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Category not found"})
	}
	database.DB.Delete(&category)
	return c.JSON(http.StatusOK, map[string]string{"message": "Category deleted"})
}
func UpdateNewsCategory(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var category models.NewsCategory
	database.DB.First(&category, id)
	if category.ID == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Category not found"})
	}
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	database.DB.Save(&category)
	return c.JSON(http.StatusOK, category)
}

func CreateNews(c echo.Context) error {
	var news models.News
	if err := c.Bind(&news); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	database.DB.Create(&news)
	return c.JSON(http.StatusCreated, news)
}

func UpdateNews(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var news models.News
	database.DB.First(&news, id)
	if news.ID == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "News not found"})
	}
	if err := c.Bind(&news); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	database.DB.Save(&news)
	return c.JSON(http.StatusOK, news)
}

func DeleteNews(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var news models.News
	database.DB.First(&news, id)
	if news.ID == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "News not found"})
	}
	database.DB.Delete(&news)
	return c.JSON(http.StatusOK, map[string]string{"message": "News deleted"})
}

func GetNews(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var news models.News
	database.DB.Preload("Category").First(&news, id)
	return c.JSON(http.StatusOK, news)
}

package handlers

import (
	"net/http"
	"strconv"

	"github.com/simon4545/goshop/database"
	"github.com/simon4545/goshop/models"

	"github.com/labstack/echo/v4"
)

func GetProductCategories(c echo.Context) error {
	var categories []models.ProductCategory
	database.DB.Find(&categories)
	return c.JSON(http.StatusOK, categories)
}

func AddProduct(c echo.Context) error {
	var product models.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err := database.DB.Create(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create product"})
	}

	return c.JSON(http.StatusCreated, product)
}

func DeleteProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
	}

	if err := database.DB.Delete(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete product"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}

func GetProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var product models.Product
	database.DB.Preload("Category").First(&product, id)
	return c.JSON(http.StatusOK, product)
}

func GetProductBySKU(c echo.Context) error {
	sku := c.Param("sku")
	var product models.Product
	if err := database.DB.Where("sku = ?", sku).Preload("Category").First(&product).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
	}
	return c.JSON(http.StatusOK, product)
}
func AddSKUs(c echo.Context) error {
	var skus []models.SKU
	if err := c.Bind(&skus); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err := database.DB.Create(&skus).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create SKUs"})
	}

	return c.JSON(http.StatusCreated, skus)
}

func DeleteSKU(c echo.Context) error {
	sku := c.Param("sku")
	var skuModel models.SKU
	if err := database.DB.Where("sku = ?", sku).First(&skuModel).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "SKU not found"})
	}

	if err := database.DB.Delete(&skuModel).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete SKU"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "SKU deleted successfully"})
}

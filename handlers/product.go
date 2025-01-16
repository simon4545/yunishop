package handlers

import (
	"net/http"
	"strconv"

	"github.com/simon4545/goshop/database"
	"github.com/simon4545/goshop/models"
	"gorm.io/gorm"

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

	return c.JSON(http.StatusOK, product)
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

func GetProductsByIDs(c echo.Context) error {
	var productIDs []uint
	if err := c.Bind(&productIDs); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	var products []models.Product
	if err := database.DB.Where("id IN ?", productIDs).Preload("Category").Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve products"})
	}

	return c.JSON(http.StatusOK, products)
}

func UpdateProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	var product models.Product
	if err := database.DB.Preload("SKUs").Preload("ProductImages").First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
	}

	var updatedProduct models.Product
	if err := c.Bind(&updatedProduct); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err := database.DB.Model(&product).Session(&gorm.Session{FullSaveAssociations: true}).Updates(updatedProduct).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update product"})
	}
	// if err1 := database.DB.Model(&product).Association("SKUs").Replace(updatedProduct.SKUs); err1 != nil {
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update product"})
	// }
	// if err2 := database.DB.Model(&product).Association("ProductImages").Replace(updatedProduct.ProductImages); err2 != nil {
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update product"})
	// }

	return c.JSON(http.StatusOK, product)
}

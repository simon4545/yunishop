package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/simon4545/goshop/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const uploadDir = "./uploads"

type Image struct {
	gorm.Model
	ProductID uint   `json:"product_id"`
	URL       string `json:"url"`
}

func UploadImage(c echo.Context) error {
	productID, err := strconv.Atoi(c.FormValue("product_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product_id"})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	// Create upload directory if it doesn't exist
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}

	// Generate unique filename
	filename := filepath.Join(uploadDir, strconv.Itoa(productID)+"_"+file.Filename)

	// Save file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Save to database
	image := models.ProductImage{
		ProductID: uint(productID),
		URL:       "/uploads/" + filepath.Base(filename),
	}
	db := c.Get("db").(*gorm.DB)
	if err := db.Create(&image).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, image)
}

func DeleteImage(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	db := c.Get("db").(*gorm.DB)
	var image models.ProductImage
	if err := db.First(&image, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "image not found"})
	}

	// Delete file
	if err := os.Remove("." + image.URL); err != nil {
		return err
	}

	// Delete from database
	if err := db.Delete(&image).Error; err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func GetImage(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	db := c.Get("db").(*gorm.DB)
	var image models.ProductImage
	if err := db.First(&image, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "image not found"})
	}

	return c.JSON(http.StatusOK, image)
}

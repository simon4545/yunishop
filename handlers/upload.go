package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

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

	return c.JSON(http.StatusOK, "/uploads/"+filepath.Base(filename))
}

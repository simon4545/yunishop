package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/exp/rand"
	"gorm.io/gorm"
)

const uploadDir = "./uploads"

type Image struct {
	gorm.Model
	ProductID uint   `json:"product_id"`
	URL       string `json:"url"`
}

func generateUniqueFilename(originalFilename string) string {
	timestamp := time.Now().Unix()
	randomNum := rand.Intn(10000)
	ext := filepath.Ext(originalFilename)
	return strconv.FormatInt(timestamp, 10) + "_" + strconv.Itoa(randomNum) + ext
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

	// Generate unique filename using the generateUniqueFilename function
	uniqueFilename := generateUniqueFilename(file.Filename)
	filename := filepath.Join(uploadDir, strconv.Itoa(productID)+"_"+uniqueFilename)

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

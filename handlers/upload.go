package handlers

import (
	"io"
	"net/http"
	"os"
	"path"
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
func ShowPic(c echo.Context) error {
	picUrl := c.Param("url")
	// 本地图片的绝对路径
	imagePath := path.Join("/mnt/kuang/project/rutian/PicBackup", picUrl)

	// 打开图片文件
	file, err := os.Open(imagePath)
	if err != nil {
		return c.String(http.StatusInternalServerError, "无法打开图片")
	}
	defer file.Close()

	// 获取图片的文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return c.String(http.StatusInternalServerError, "无法获取图片信息")
	}

	// 设置响应头
	c.Response().Header().Set(echo.HeaderContentType, "image/jpeg")
	c.Response().Header().Set(echo.HeaderContentLength, string(fileInfo.Size()))

	// 将图片内容写入响应
	return c.Stream(http.StatusOK, "image/jpeg", file)
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

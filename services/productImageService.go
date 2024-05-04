package services

import (
	"etsy-spider/models"
	"fmt"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ProductImageService struct {
	Db *gorm.DB
}

func (s ProductImageService) AddProductImage(productId uint, imageUrl, sku, imageDir string) (result *gorm.DB) {
	_, filePath := s.downloadImage(imageUrl, sku, imageDir)
	data := models.ProductImage{ProductId: productId, Image: imageUrl, SavePath: filePath}
	result = s.Db.Create(&data)
	return result
}

// 从URL中提取文件名
func filenameFromURL(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

// 根据URL下载图片到指定的文件
func (s ProductImageService) downloadImage(url, sku, imageDir string) (err error, filePath string) {
	// 获取图片文件名
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	// 创建 SKU 目录
	dir := filepath.Join(imageDir, sku)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err, ""
	}

	filePath = filepath.Join(dir, fileName)
	_, err = os.Stat(filePath)
	if err == nil {
		return nil, filePath
	}
	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return err, ""
	}
	defer file.Close()

	err = s.wget(url, file)
	if err != nil {
		return err, ""
	}
	return nil, filePath
}

func (s ProductImageService) wget(url string, file *os.File) (err error) {
	// 发起 HTTP GET 请求下载图片
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// 检查响应状态码
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image, status code: %d", response.StatusCode)
	}

	// 将响应体内容保存到文件
	_, err = io.Copy(file, response.Body)
	return err
}

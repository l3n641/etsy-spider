package services

import (
	"etsy-spider/models"
	"gorm.io/gorm"
)

type ProductService struct {
	Db *gorm.DB
}

func (s ProductService) AddProduct(url, sku, title, html string) (result *gorm.DB, id uint) {
	data := models.Product{Url: url, Sku: sku, Title: title, HtmlContent: html}
	result = s.Db.Create(&data)
	return result, data.ID
}

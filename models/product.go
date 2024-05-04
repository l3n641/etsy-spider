package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Title       string
	Sku         string
	Url         string `gorm:"index"`
	HtmlContent string
}

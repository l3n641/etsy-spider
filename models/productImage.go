package models

type ProductImage struct {
	ProductId uint `gorm:"index"`
	Image     string
	SavePath  string
}

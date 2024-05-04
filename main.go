package main

import (
	"crypto/md5"
	"encoding/hex"
	"etsy-spider/api"
	"etsy-spider/models"
	"etsy-spider/services"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net/http"
	"os"
)

type Configuration struct {
	WebServiceAddr string `mapstructure:"webServiceAddr" json:"webServiceAddr" yaml:"webServiceAddr"`
	DbName         string `mapstructure:"dbName" json:"dbName" yaml:"dbName"`
	ImageDir       string `mapstructure:"imageDir" json:"imageDir" yaml:"imageDir"`
}

func initializeConfig() *Configuration {
	var config Configuration

	// 设置配置文件路径
	configFile := "./conf.yaml"
	// 生产环境可以通过设置环境变量来改变配置文件路径
	if configEnv := os.Getenv("VIPER_CONFIG"); configEnv != "" {
		configFile = configEnv
	}

	// 初始化 viper
	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read configFile failed: %s \n", err))
	}

	// 设置默认值
	viper.SetDefault("webServiceAddr", ":8080")
	viper.SetDefault("imageDir", "./static/image")

	if err := v.Unmarshal(&config); err != nil {
		fmt.Println(err)
	}

	return &config
}

func __init__() {
}

func main() {
	var AppConfig = initializeConfig()

	db, err := gorm.Open(sqlite.Open(AppConfig.DbName), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Product{}, &models.ProductImage{})

	r := gin.Default()
	r.StaticFS("/image", http.Dir(AppConfig.ImageDir))
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	r.Use(cors.New(config))

	r.POST("/product", func(c *gin.Context) {
		var data api.ProductReq
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		productSrc := services.ProductService{Db: db}
		sum := md5.Sum([]byte(data.Url))
		sku := hex.EncodeToString(sum[:])
		result, productId := productSrc.AddProduct(data.Url, sku, data.Title, data.HtmlContent)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "保存产品失败",
			})
			return
		}

		go func() {
			imageSrc := services.ProductImageService{Db: db}
			for _, image := range data.Images {
				imageSrc.AddProductImage(productId, image, sku, AppConfig.ImageDir)
			}
		}()

		c.JSON(http.StatusOK, gin.H{
			"id": productId,
		})
		return

	})

	r.Run(AppConfig.WebServiceAddr)
}

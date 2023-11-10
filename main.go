package main

import (
	"go-qrcode-generator-cms-api/docs"
	"go-qrcode-generator-cms-api/src/configs"
	"go-qrcode-generator-cms-api/src/entity"
	"go-qrcode-generator-cms-api/src/transport/rest"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Go QRCode Generator CMS - Swagger API Discovery
//	@version		1.0
//	@description	Go QRCode Generator CMS - Swagger API Discovery

//	@contact.name	Xuan Hoa Le
//	@contact.email	dso.intern.xuanhoa@gmail.com

// @host		localhost:3000
// @BasePath	/api/v1
func main() {
	if db, err := configs.NewGormClient().Instance(); err != nil {
		panic("Can't connect to database via GORM: " + err.Error())
	} else if redisClient, err := configs.NewRedisCacheClient().Instance(); err != nil {
		panic("Can't connect to Redis Server via Redis Client: " + err.Error())
	} else if cld, err := configs.NewCloudinaryClient().Instance(); err != nil {
		panic("Can't connect to Cloudinary via Cloudinary API: " + err.Error())
	} else if oauth2cfg, err := configs.NewOAuthClient().Instance(); err != nil {
		panic("Can't connect to Google Authentication Service via OAuth2: " + err.Error())
	} else {
		port := os.Getenv("PORT")
		models := []interface{}{
			&entity.Role{},
			&entity.User{},
			&entity.QRCode{},
		}
		db.AutoMigrate(models...)

		currentDir, _ := os.Getwd()
		viewsDir := filepath.Join(currentDir, "./static/views/*")

		router := gin.Default()
		router.LoadHTMLGlob(viewsDir)
		rest.NewRouteConfig(router).Config(db, redisClient, cld, oauth2cfg)

		docs.SwaggerInfo.BasePath = "/api/v1"
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		router.Run(":" + port)
	}
}

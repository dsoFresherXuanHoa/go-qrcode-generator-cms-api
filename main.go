package main

import (
	"go-qrcode-generator-cms-api/src/configs"
	"go-qrcode-generator-cms-api/src/entity"
	"go-qrcode-generator-cms-api/src/transport/rest"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	if db, err := configs.NewGormInstance().GetGormInstance(); err != nil {
		panic("Can't connect to database via GORM: " + err.Error())
	} else if cld, err := configs.NewCloudinaryInstance().GetCloudinaryInstance(); err != nil {
		panic("Can't connect to cloudinary server via Cloudinary API: " + err.Error())
	} else if oauth2cfg, err := configs.NewOAuthInstance().GetOAuthConfigInstance(); err != nil {
		panic("Can't connect to Google authentication server via OAuth2: " + err.Error())
	} else {
		port := os.Getenv("PORT")
		models := []interface{}{
			&entity.Role{},
			&entity.User{},
		}
		db.AutoMigrate(models...)

		currentDir, _ := os.Getwd()
		viewsDir := filepath.Join(currentDir, "./static/views/*")
		router := gin.Default()
		router.LoadHTMLGlob(viewsDir)
		rest.NewRouteConfig(router).RouteConfig(db, cld, oauth2cfg)
		router.Run(":" + port)
	}
}

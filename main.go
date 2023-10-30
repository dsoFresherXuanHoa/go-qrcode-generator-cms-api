package main

import (
	"go-qrcode-generator-cms-api/src/configs"
	"go-qrcode-generator-cms-api/src/entity"
	"go-qrcode-generator-cms-api/src/transport/rest"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	if db, err := configs.NewGormInstance().GetGormInstance(); err != nil {
		panic("Can't connect to database via GORM: " + err.Error())
	} else if cld, err := configs.NewCloudinaryInstance().GetCloudinaryInstance(); err != nil {
		panic("Can't connect to cloudinary server via Cloudinary API: " + err.Error())
	} else {
		port := os.Getenv("PORT")
		models := []interface{}{
			&entity.Role{},
			&entity.User{},
		}
		db.AutoMigrate(models...)

		router := gin.Default()
		rest.NewRouteConfig(router).RouteConfig(db, cld)
		router.Run(":" + port)
	}
}

package rest

import (
	"go-qrcode-generator-cms-api/src/middlewares"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type routeConfig struct {
	router *gin.Engine
}

func NewRouteConfig(router *gin.Engine) *routeConfig {
	return &routeConfig{router: router}
}

func (cfg routeConfig) RouteConfig(db *gorm.DB, cld *cloudinary.Cloudinary) {
	secretKey := os.Getenv("JWT_ACCESS_SECRET")
	v1 := cfg.router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/sign-up", SignUp(db, cld))
			auth.GET("/activation", Activate(db))
			auth.POST("/sign-in", SignIn(db))
			auth.GET("/me", middlewares.RequiredAuthorized(db, secretKey), Me(db))
		}

		role := v1.Group("/role")
		{
			role.POST("/", CreateRole(db))
		}
	}
}

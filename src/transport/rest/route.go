package rest

import (
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
	v1 := cfg.router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/sign-up", SignUp(db, cld))
			auth.GET("/activation", Activate(db))
		}

		role := v1.Group("/role")
		{
			role.POST("/", CreateRole(db))
		}
	}
}

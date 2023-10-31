package rest

import (
	"go-qrcode-generator-cms-api/src/middlewares"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type routeConfig struct {
	router *gin.Engine
}

func NewRouteConfig(router *gin.Engine) *routeConfig {
	return &routeConfig{router: router}
}

func (cfg routeConfig) RouteConfig(db *gorm.DB, cld *cloudinary.Cloudinary, oauth2cfg *oauth2.Config) {
	secretKey := os.Getenv("JWT_ACCESS_SECRET")
	v1 := cfg.router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.GET("/", Home(db))
			auth.POST("/sign-up", SignUp(db, cld))
			auth.PATCH("/activation", Activate(db))
			auth.POST("/sign-in", SignIn(db))
			auth.GET("/me", middlewares.RequiredAuthorized(db, secretKey), Me(db))
			auth.GET("/reset-password", RequestResetPassword(db))
			auth.PATCH("/reset-password", ResetPassword(db))

			oauth := auth.Group("/oauth")
			{
				oauth.GET("/sign-in", GoogleSignIn(db, oauth2cfg))
				oauth.GET("/callback", GoogleSignInCallBack(db, oauth2cfg))
			}
		}

		role := v1.Group("/role")
		{
			role.POST("/", CreateRole(db))
		}
	}
}

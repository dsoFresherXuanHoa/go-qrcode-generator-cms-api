package rest

import (
	"go-qrcode-generator-cms-api/src/middlewares"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type routeConfig struct {
	router *gin.Engine
}

func NewRouteConfig(router *gin.Engine) *routeConfig {
	return &routeConfig{router: router}
}

func (cfg routeConfig) Config(db *gorm.DB, redisClient *redis.Client, cld *cloudinary.Cloudinary, oauth2cfg *oauth2.Config) {
	secretKey := os.Getenv("JWT_ACCESS_SECRET")
	cfg.router.MaxMultipartMemory = 8 << 20
	v1 := cfg.router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.GET("/", Home(db))
			auth.POST("/sign-up", SignUp(db, cld))
			auth.PATCH("/activation", Activate(db))
			auth.POST("/sign-in", SignIn(db, redisClient))
			auth.GET("/me", middlewares.RequiredAuthorized(db, redisClient, secretKey), Me(db))
			auth.GET("/reset-password", RequestResetPassword(db))
			auth.PATCH("/reset-password", ResetPassword(db, redisClient))
			auth.GET("/sign-out", middlewares.RequiredAuthorized(db, redisClient, secretKey), SignOut(redisClient))

			oauth := auth.Group("/oauth")
			{
				oauth.GET("/sign-in", GoogleSignIn(db, oauth2cfg))
				oauth.GET("/callback", GoogleSignInCallBack(db, redisClient, oauth2cfg))
			}
		}

		roles := v1.Group("/roles")
		{
			roles.POST("/", CreateRole(db))
		}

		users := v1.Group("/users")
		{
			users.GET("/:userUUID/qrcodes/", middlewares.RequiredAdministratorPermission(db, redisClient, secretKey), FindQRCodeByUserId(db))
		}

		qrcodes := v1.Group("/qrcodes")
		{
			qrcodes.POST("/", middlewares.RateLimit(), middlewares.RequiredAuthorized(db, redisClient, secretKey), CreateQRCode(db, redisClient, cld))
			qrcodes.GET("/:uuid", middlewares.RequiredAdministratorPermission(db, redisClient, secretKey), FindQRCodeByUUID(db))
			qrcodes.GET("/", middlewares.RequiredAdministratorPermission(db, redisClient, secretKey), FindQRCodeByCondition(db))
		}
	}
}

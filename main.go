package main

import (
	"fmt"
	"go-qrcode-generator-cms-api/docs"
	"go-qrcode-generator-cms-api/proto/gRPC/qrcodes"
	"go-qrcode-generator-cms-api/src/configs"
	"go-qrcode-generator-cms-api/src/entity"
	"go-qrcode-generator-cms-api/src/transport/rest"
	"go-qrcode-generator-cms-api/src/transport/rpc"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gammazero/workerpool"
	"github.com/gin-gonic/gin"

	"google.golang.org/grpc"

	cors "github.com/rs/cors/wrapper/gin"
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
		// Restful Service
		port := os.Getenv("PORT")
		wpSize, _ := strconv.Atoi(os.Getenv("WORKER_POOL_SIZE"))
		wp := workerpool.New(wpSize)
		models := []interface{}{
			&entity.Role{},
			&entity.User{},
			&entity.QRCode{},
		}
		db.AutoMigrate(models...)

		currentDir, _ := os.Getwd()
		viewsDir := filepath.Join(currentDir, "./static/views/*")

		router := gin.Default()
		router.Use(cors.AllowAll())
		router.LoadHTMLGlob(viewsDir)

		rest.NewRouteConfig(router).Config(wp, db, redisClient, cld, oauth2cfg)

		docs.SwaggerInfo.BasePath = "/api/v1"
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		go router.Run(":" + port)

		// GRPC Service
		grpcServiceAddress := os.Getenv("GRPC_SERVICE_ADDRESS")
		if lis, err := net.Listen("tcp", grpcServiceAddress); err != nil {
			panic("Error while start gRPC server (qrcode service) at: " + grpcServiceAddress + " with error: " + err.Error())
		} else {
			fmt.Println("QRCode Service is running at: " + grpcServiceAddress)
			s := grpc.NewServer()
			server := rpc.NewGrpcServer(db, redisClient)
			qrcodes.RegisterQRCodeServiceServer(s, server)
			s.Serve(lis)
		}
	}
}

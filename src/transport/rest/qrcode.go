package rest

import (
	"errors"
	"fmt"
	"go-qrcode-generator-cms-api/src/business"
	"go-qrcode-generator-cms-api/src/constants"
	"go-qrcode-generator-cms-api/src/entity"
	"go-qrcode-generator-cms-api/src/storage"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	InvalidCreateQRCodeRequest = "Invalid QR Code Incoming Request: Check Swagger For More Information."
	QRCodeUUIDEmpty            = "Invalid QR Code Incoming Request: QRCode UUID can not be empty."

	CreateQrCodeFailure    = "Cannot Create QR Code: Make Sure You Has Right Permission And Try Again."
	GetQRCodeByUUIDFailure = "Cannot Get QR Code By UUID: Make Sure You Has Right Permission And Try Again."

	CreateQrCodeSuccess    = "Create QR Code Success: Congrats."
	GetQRCodeByUUIDSuccess = "Get QR Code By UUID Success: Congrats."

	ErrQRCodeUUIDEmpty = errors.New("get qrcode uuid from user request")
)

func CreateQRCode(db *gorm.DB, redisClient *redis.Client, cld *cloudinary.Cloudinary) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sqlStorage := storage.NewSQLStore(db)
		redisStorage := storage.NewRedisStore(redisClient)
		qrCodeStorage := storage.NewQrCodeStore(sqlStorage, redisStorage)
		qrCodeBusiness := business.NewQRCodeBusiness(qrCodeStorage, redisStorage)

		var reqQrCode entity.QRCodeCreatable
		if err := ctx.ShouldBind(&reqQrCode); err != nil {
			fmt.Println("Error while parse user request: " + err.Error())
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, ErrQRCodeUUIDEmpty.Error(), QRCodeUUIDEmpty))
		} else {
			userId := ctx.Value("userId").(uint)
			reqQrCode.UserID = userId
			if _, publicURL, err := qrCodeBusiness.CreateQRCode(ctx, redisClient, cld, &reqQrCode); err != nil {
				fmt.Println("Error while create QrCode: " + err.Error())
				ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), CreateQrCodeFailure))
			} else {
				ctx.JSON(http.StatusOK, entity.NewStandardResponse(gin.H{"publicURL": publicURL, "encode": nil}, http.StatusOK, "OK", "", CreateQrCodeSuccess))
			}
		}
	}
}

func FindQRCodeByUUID(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sqlStorage := storage.NewSQLStore(db)
		qrCodeStorage := storage.NewQrCodeStore(sqlStorage, nil)
		qrCodeBusiness := business.NewQRCodeBusiness(qrCodeStorage, nil)

		if qrCodeUUID := ctx.Param("uuid"); qrCodeUUID == "" {
			fmt.Println("Error while get qrcode uuid from user request: qrcode uuid can not be empty")
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, ErrQRCodeUUIDEmpty.Error(), QRCodeUUIDEmpty))
		} else if qrCode, err := qrCodeBusiness.FindQRCodeByUUID(ctx, qrCodeUUID); err != nil {
			fmt.Println("Error while get qrcode by uuid: " + err.Error())
			ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), GetQRCodeByUUIDFailure))
		} else {
			ctx.JSON(http.StatusOK, entity.NewStandardResponse(qrCode, http.StatusOK, "OK", "", GetQRCodeByUUIDSuccess))
		}
	}
}

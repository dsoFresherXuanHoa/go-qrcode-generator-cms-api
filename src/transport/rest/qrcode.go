package rest

import (
	"fmt"
	"go-qrcode-generator-cms-api/src/business"
	"go-qrcode-generator-cms-api/src/constants"
	"go-qrcode-generator-cms-api/src/entity"
	"go-qrcode-generator-cms-api/src/storage"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	InvalidCreateQRCodeRequest = "Invalid QR Code Incoming Request: Check Swagger For More Information."
	CreateQrCodeFailure        = "Cannot Create QR Code: Make Sure You Has Right Permission And Try Again."
	CreateQrCodeSuccess        = "Create QR Code Success: Congrats."
)

func CreateQRCode(db *gorm.DB, cld *cloudinary.Cloudinary) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sqlStorage := storage.NewSQLStore(db)
		qrCodeStorage := storage.NewQrCodeStore(sqlStorage)
		qrCodeBusiness := business.NewQRCodeBusiness(qrCodeStorage)

		var reqQrCode entity.QRCodeCreatable
		if err := ctx.ShouldBind(&reqQrCode); err != nil {
			fmt.Println("Error while parse user request: " + err.Error())
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, err.Error(), InvalidCreateQRCodeRequest))
		} else {
			userId := ctx.Value("userId").(uint)
			reqQrCode.UserID = userId
			if encode, publicURL, err := qrCodeBusiness.CreateQRCode(ctx, cld, &reqQrCode); err != nil {
				fmt.Println("Error while create QrCode: " + err.Error())
				ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), CreateQrCodeFailure))
			} else {
				ctx.JSON(http.StatusOK, entity.NewStandardResponse(gin.H{"publicURL": publicURL, "encode": encode}, http.StatusOK, "OK", "", CreateQrCodeSuccess))
			}
		}
	}
}

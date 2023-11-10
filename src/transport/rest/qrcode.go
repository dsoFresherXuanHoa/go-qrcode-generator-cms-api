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
	InvalidPaginateRequest     = "Invalid Paginate Incoming Request: paging must include page and size."

	CreateQrCodeFailure                = "Cannot Create QR Code: Make Sure You Has Right Permission And Try Again."
	GetQRCodeByUUIDFailure             = "Cannot Get QR Code By UUID: Make Sure You Has Right Permission And Try Again."
	GetQRCodeByConditionFailure        = "Cannot Get QR Code By Condition: Make Sure You Has Right Permission And Try Again."
	GetAllQRCodeFailure                = "Cannot Get All QR Code: Make Sure You Has Right Permission And Try Again."
	ValidateCreateQRCodeRequestFailure = "Invalid Create QRCode Incoming Request: Check Swagger For More Information."
	QRCodeNotFound                     = "Cannot Get QR Code By UUID: QR Code Not Found."

	CreateQrCodeSuccess         = "Create QR Code Success: Congrats."
	GetQRCodeByUUIDSuccess      = "Get QR Code By UUID Success: Congrats."
	GetQRCodeByConditionSuccess = "Get QR Code By Condition Success: Congrats."
	GetAllQRCodeSuccess         = "Get QR All Code Success: Congrats."

	ErrQRCodeUUIDEmpty = errors.New("get qrcode uuid from user request")
)

// CreateQRCode godoc
//
//	@Summary		Create QRCode
//	@Description	Create QRCode using custom configuration
//	@Tags			qrcodes
//	@Accept       multipart/form-data
//	@Produce		json
//	@Param  Authorization  header  string  required  "Bearer Token"
//	@Param		 qrcode	formData	entity.QRCodeCreatable	true	"QRCode"
//	@Param		 logo	formData	file	false	"Logo"
//	@Param		 halftone	formData	file	false	"Halftone"
//	@Success		200		{object}	entity.standardResponse
//	@Failure		400		{object}	entity.standardResponse
//	@Failure		500		{object}	entity.standardResponse
//	@Router			/qrcodes [post]
func CreateQRCode(db *gorm.DB, redisClient *redis.Client, cld *cloudinary.Cloudinary) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sqlStorage := storage.NewSQLStore(db)
		redisStorage := storage.NewRedisStore(redisClient)
		qrCodeStorage := storage.NewQrCodeStore(sqlStorage, redisStorage)
		qrCodeBusiness := business.NewQRCodeBusiness(qrCodeStorage, redisStorage)

		var reqQrCode entity.QRCodeCreatable
		if err := ctx.ShouldBind(&reqQrCode); err != nil {
			fmt.Println("Error while parse user request: " + err.Error())
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, err.Error(), InvalidCreateQRCodeRequest))
		} else if err := reqQrCode.Validate(); err != nil {
			fmt.Println("Error while validate user request: " + err.Error())
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, err.Error(), ValidateCreateQRCodeRequestFailure))
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

// FindQRCodeByUUID godoc
//
//	@Summary		Find QRCode by UUID
//	@Description	Find QRCode using UUID - full QRCode Information
//	@Tags			qrcodes
//	@Accept			json
//	@Produce		json
//	@Param  Authorization  header  string  required  "Bearer Token"
//	@Param			qrCodeUUID	path		string	true	"QRCode UUID"
//	@Success		200	{object}	entity.standardResponse
//	@Failure		400	{object}	entity.standardResponse
//	@Failure		404	{object}	entity.standardResponse
//	@Failure		500	{object}	entity.standardResponse
//	@Router			/qrcodes/{qrCodeUUID} [get]
func FindQRCodeByUUID(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sqlStorage := storage.NewSQLStore(db)
		qrCodeStorage := storage.NewQrCodeStore(sqlStorage, nil)
		qrCodeBusiness := business.NewQRCodeBusiness(qrCodeStorage, nil)

		if qrCodeUUID := ctx.Param("uuid"); qrCodeUUID == "" {
			fmt.Println("Error while get qrcode uuid from user request: qrcode uuid can not be empty")
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, ErrQRCodeUUIDEmpty.Error(), QRCodeUUIDEmpty))
		} else if qrCode, err := qrCodeBusiness.FindQRCodeByUUID(ctx, qrCodeUUID); err == storage.ErrFindQRCodeByUUID {
			fmt.Println("Error while get qrcode by uuid: " + err.Error())
			ctx.JSON(http.StatusNotFound, entity.NewStandardResponse(nil, http.StatusNotFound, constants.StatusNotFound, err.Error(), QRCodeNotFound))
		} else if err != nil {
			fmt.Println("Error while get qrcode by uuid: " + err.Error())
			ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), GetQRCodeByUUIDFailure))
		} else {
			ctx.JSON(http.StatusOK, entity.NewStandardResponse(qrCode, http.StatusOK, "OK", "", GetQRCodeByUUIDSuccess))
		}
	}
}

// FindQRCodeByCondition godoc
//
//	@Summary		Find QRCode by custom condition
//	@Description	Find QRCode by custom condition
//	@Tags			qrcodes
//	@Accept			json
//	@Produce		json
//	@Param  Authorization  header  string  required  "Bearer Token"
//	@Param			page	query		int	false	"Page"
//	@Param			size	query		int	false	"Size"
//	@Param			version	query		int	false	"Version"
//	@Param			type	query		string	false	"Type"
//	@Param			errorLevel	query		int	false	"Error Level"
//	@Param			startTime	query		int	false	"Start Time" config:"format=rfc3339"
//	@Param			endTime	query		int	false	"End Time" config:"format=rfc3339"
//	@Success		200	{object}		entity.standardResponse
//	@Failure		400	{object}	entity.standardResponse
//	@Failure		500	{object}	entity.standardResponse
//	@Router			/qrcodes [get]
func FindQRCodeByCondition(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sqlStorage := storage.NewSQLStore(db)
		qrCodeStorage := storage.NewQrCodeStore(sqlStorage, nil)
		qrCodeBusiness := business.NewQRCodeBusiness(qrCodeStorage, nil)

		cond := map[string]interface{}{}
		timeStat := map[string]string{}
		if qrCodeType := ctx.Query("type"); qrCodeType != "" {
			cond["type"] = qrCodeType
		}
		if qrCodeVersion := ctx.Query("version"); qrCodeVersion != "" {
			cond["version"] = qrCodeVersion
		}
		if qrCodeErrorLevel := ctx.Query("errorLevel"); qrCodeErrorLevel != "" {
			cond["error_level"] = qrCodeErrorLevel
		}
		if startTime := ctx.Query("startTime"); startTime != "" {
			timeStat["start_time"] = startTime
			if endTime := ctx.Query("endTime"); endTime != "" {
				timeStat["end_time"] = endTime
			}
		}

		var paging entity.Paginate
		if err := ctx.ShouldBind(&paging); err != nil {
			fmt.Println("Error while parse paging from user request: " + err.Error())
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, err.Error(), GetQRCodeByConditionFailure))
		} else if qrCodes, err := qrCodeBusiness.FindQRCodeByCondition(ctx, cond, timeStat, &paging); err != nil {
			fmt.Println("Error while find all qrcode by condition: " + err.Error())
			ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), GetQRCodeByConditionFailure))
		} else {
			ctx.JSON(http.StatusOK, entity.NewStandardWithPaginateResponse(qrCodes, http.StatusOK, "OK", "", GetQRCodeByConditionSuccess, paging))
		}
	}
}

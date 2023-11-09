package rest

import (
	"fmt"
	"go-qrcode-generator-cms-api/src/business"
	"go-qrcode-generator-cms-api/src/constants"
	"go-qrcode-generator-cms-api/src/entity"
	"go-qrcode-generator-cms-api/src/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	InvalidUserId = "Invalid User ID Incoming Request: User ID can not be empty and must be numeric."

	GetQRCodeByUserIdFailure = "Cannot Get All QR Code By User Id: Make Sure You Has Right Permission And Try Again."

	GetQRCodeByUserIdSuccess = "Get All QR Code By User Id Success: Congrats."
)

func FindQRCodeByUserId(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sqlStorage := storage.NewSQLStore(db)
		userStorage := storage.NewUserStore(sqlStorage)
		userBusiness := business.NewUserBusiness(userStorage)

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
		if userId, err := strconv.Atoi(ctx.Param("userId")); err != nil {
			fmt.Println("Error while get userId from user request: " + err.Error())
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, err.Error(), InvalidUserId))
		} else if err := ctx.ShouldBind(&paging); err != nil {
			fmt.Println("Error while parse paging from user request: " + err.Error())
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, err.Error(), InvalidPaginateRequest))
		} else {
			id := uint(userId)
			if qrCodes, err := userBusiness.FindQRCodeByUserId(ctx, id, cond, timeStat, paging); err != nil {
				fmt.Println("Error while find all qrcode by userId: " + err.Error())
				ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), GetQRCodeByUserIdFailure))
			} else {
				ctx.JSON(http.StatusOK, entity.NewStandardWithPaginateResponse(qrCodes, http.StatusOK, "OK", "", GetQRCodeByUserIdSuccess, paging))
			}
		}
	}
}

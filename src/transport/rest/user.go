package rest

import (
	"errors"
	"fmt"
	"go-qrcode-generator-cms-api/src/business"
	"go-qrcode-generator-cms-api/src/constants"
	"go-qrcode-generator-cms-api/src/entity"
	"go-qrcode-generator-cms-api/src/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	ErrUserUUIDEmpty = errors.New("get user uuid from user request")

	UserUUIDEmpty = "Invalid User ID Incoming Request: User ID can not be empty and must be uuid format."

	GetQRCodeByUserIdFailure = "Cannot Get All QR Code By User Id: Make Sure You Has Right Permission And Try Again."

	GetQRCodeByUserIdSuccess = "Get All QR Code By User Id Success: Congrats."
)

// FindQRCodeByUserId godoc
//
//	@Summary		Find QRCode by user id
//	@Description	Find QRCode by user id with custom condition
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param  Authorization  header  string  required  "Bearer Token"
//	@Param			userUUID	path		string	true	"User UUID"
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
//	@Router			/users/{userUUID}/qrcodes [get]
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
		if userUUID := ctx.Param("userUUID"); userUUID == "" {
			fmt.Println("Error while get userUUID from user request: empty userUUID")
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, ErrUserUUIDEmpty.Error(), UserUUIDEmpty))
		} else if err := ctx.ShouldBind(&paging); err != nil {
			fmt.Println("Error while parse paging from user request: " + err.Error())
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, err.Error(), InvalidPaginateRequest))
		} else {
			cond["user_id"] = userUUID
			if qrCodes, err := userBusiness.FindQRCodeByUserId(ctx, userUUID, cond, timeStat, &paging); err != nil {
				fmt.Println("Error while find all qrcode by userId: " + err.Error())
				ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), GetQRCodeByUserIdFailure))
			} else {
				ctx.JSON(http.StatusOK, entity.NewStandardWithPaginateResponse(qrCodes, http.StatusOK, "OK", "", GetQRCodeByUserIdSuccess, paging))
			}
		}
	}
}

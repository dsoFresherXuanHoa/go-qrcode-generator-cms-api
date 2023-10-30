package rest

import (
	"fmt"
	"go-qrcode-generator-cms-api/src/business"
	"go-qrcode-generator-cms-api/src/constants"
	"go-qrcode-generator-cms-api/src/entity"
	"go-qrcode-generator-cms-api/src/storage"
	"go-qrcode-generator-cms-api/src/utils"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SignUp(db *gorm.DB, cld *cloudinary.Cloudinary) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cloudinaryStorage := storage.NewCloudinaryStore(cld)
		sqlStorage := storage.NewSQLStore(db)
		userStorage := storage.NewUserStore(sqlStorage)
		authStorage := storage.NewAuthStore(userStorage)
		authBusiness := business.NewAuthBusiness(authStorage)

		var reqUser entity.UserCreatable
		if err := ctx.ShouldBind(&reqUser); err != nil {
			fmt.Println("Error while parse user request to struct: " + err.Error())
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, err.Error(), constants.InvalidSignUpRequestFormat))
		} else {
			if encodeAvatar, err := utils.NewImageUtil().ImageFileHeader2Base64(reqUser.Avatar); err != nil {
				fmt.Println("Error while encode file multipart header to base64 format: " + err.Error())
				ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), constants.ErrEncodeFileMultiPartHeader))
			} else if uploadResult, err := cloudinaryStorage.UploadSingleImage(ctx, *encodeAvatar); err != nil {
				fmt.Println("Error while upload single image to Cloudinary API: " + err.Error())
				ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), constants.ErrUploadSingleFileToCloudinary))
			} else {
				reqUser.AvatarURL = uploadResult.URL
				if userUUID, err := authBusiness.SignUp(ctx, &reqUser); err != nil {
					fmt.Println("Error while sign up for new user in auth transport: " + err.Error())
					ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), constants.ErrSignUpForNewUser))
				} else if err := utils.NewMailUtil().SendActivationRequestEmail(reqUser); err != nil {
					fmt.Println("Error while send activation request email to user in auth transport: " + err.Error())
					ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), constants.ErrSendActivationRequestEmail))
				} else {
					ctx.JSON(http.StatusOK, entity.NewStandardResponse(userUUID, http.StatusOK, "OK", "", constants.SignUpForNewUserSuccess))
				}
			}
		}
	}
}

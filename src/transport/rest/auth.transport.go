package rest

import (
	"fmt"
	"go-qrcode-generator-cms-api/src/business"
	"go-qrcode-generator-cms-api/src/constants"
	"go-qrcode-generator-cms-api/src/entity"
	"go-qrcode-generator-cms-api/src/errors"
	"go-qrcode-generator-cms-api/src/storage"
	"go-qrcode-generator-cms-api/src/utils"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
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

func Activate(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sqlStorage := storage.NewSQLStore(db)
		userStorage := storage.NewUserStore(sqlStorage)
		authStorage := storage.NewAuthStore(userStorage)
		authBusiness := business.NewAuthBusiness(authStorage)

		if activationCode := ctx.Query("activationCode"); activationCode == "" {
			fmt.Println("Error while get activationCode from user request in auth transport: missing activation code from query string")
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, errors.ErrMissingActivationCodeInQueryString.Error(), constants.MissingActivationCodeInQueryString))
		} else if err := authBusiness.Activate(ctx, activationCode); err != nil {
			fmt.Println("Error while activate user in auth transport: " + err.Error())
			ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), constants.ErrActivateUser))
		} else {
			ctx.JSON(http.StatusOK, entity.NewStandardResponse(true, http.StatusOK, constants.StatusOK, "", constants.ActivateUserSuccess))
		}
	}
}

func SignIn(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sqlStorage := storage.NewSQLStore(db)
		userStorage := storage.NewUserStore(sqlStorage)
		authStorage := storage.NewAuthStore(userStorage)
		authBusiness := business.NewAuthBusiness(authStorage)

		var reqUser entity.UserQueryable
		if err := ctx.ShouldBind(&reqUser); err != nil {
			fmt.Println("Error while parse user request to struct: " + err.Error())
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, err.Error(), constants.InvalidUserQueryRequestFormat))
		} else {
			if accessToken, err := authBusiness.SignIn(ctx, &reqUser); err != nil {
				fmt.Println("Error while sign in in auth transport: " + err.Error())
				ctx.JSON(http.StatusUnauthorized, entity.NewStandardResponse(nil, http.StatusUnauthorized, constants.StatusUnauthorized, err.Error(), constants.ErrSignIn))
			} else {
				ctx.JSON(http.StatusOK, entity.NewStandardResponse(gin.H{"accessToken": accessToken}, http.StatusOK, constants.StatusOK, "", constants.SignInSuccess))
			}
		}
	}
}

func Me(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sqlStorage := storage.NewSQLStore(db)
		userStorage := storage.NewUserStore(sqlStorage)
		authStorage := storage.NewAuthStore(userStorage)
		authBusiness := business.NewAuthBusiness(authStorage)

		userId := ctx.Value("userId").(uint)
		if user, err := authBusiness.Me(ctx, userId); err != nil {
			fmt.Println("Error while find detail user in auth transport: " + err.Error())
			ctx.JSON(http.StatusForbidden, entity.NewStandardResponse(nil, http.StatusForbidden, constants.StatusForbidden, err.Error(), constants.ErrFindDetailUser))
		} else {
			ctx.JSON(http.StatusOK, entity.NewStandardResponse(user, http.StatusOK, constants.StatusOK, "", constants.FindDetailUserSuccess))
		}
	}
}

func RequestResetPassword(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sqlStorage := storage.NewSQLStore(db)
		userStorage := storage.NewUserStore(sqlStorage)
		if email := ctx.Query("email"); email == "" {
			fmt.Println("Error while get email from user request in auth transport: missing email from query string")
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, errors.ErrMissingEmailInQueryString.Error(), constants.MissingEmailInQueryString))
		} else if usr, err := userStorage.FindUserByEmail(ctx, email); err != nil {
			fmt.Println("Error while find user by email in auth transport: " + err.Error())
			ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), constants.ErrUserNotFound))
		} else if err := utils.NewMailUtil().SendResetPasswordRequestEmail(*usr); err != nil {
			fmt.Println("Error while send reset password email reset to user email: " + err.Error())
			ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), constants.ErrSendResetPasswordRequestEmail))
		} else {
			ctx.JSON(http.StatusOK, entity.NewStandardResponse(true, http.StatusOK, constants.StatusOK, "", constants.SendResetPasswordRequestEmailSuccess))
		}
	}
}

func ResetPassword(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sqlStorage := storage.NewSQLStore(db)
		userStorage := storage.NewUserStore(sqlStorage)
		authStorage := storage.NewAuthStore(userStorage)
		authBusiness := business.NewAuthBusiness(authStorage)

		var reqUser entity.UserUpdatable
		if err := ctx.ShouldBind(&reqUser); err != nil {
			fmt.Println("Error while parse user request to struct: " + err.Error())
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, err.Error(), constants.InvalidUserQueryRequestFormat))
		} else if resetCode := ctx.Query("resetCode"); resetCode == "" {
			fmt.Println("Error while get activationCode from user request in auth transport: missing activation code from query string")
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, errors.ErrMissingActivationCodeInQueryString.Error(), constants.MissingActivationCodeInQueryString))
		} else if err := authBusiness.ResetPassword(ctx, resetCode, &reqUser); err != nil {
			fmt.Println("Error while reset user password in auth transport: " + err.Error())
			ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), constants.ErrResetPassword))
		} else {
			ctx.JSON(http.StatusOK, entity.NewStandardResponse(true, http.StatusOK, constants.StatusOK, "", constants.ResetPasswordSuccess))
		}
	}
}

func Home(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "sign-in.html", nil)
	}
}

func GoogleSignIn(db *gorm.DB, oauth2 *oauth2.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: Random State to avoid CSRF attack:
		url := oauth2.AuthCodeURL("dsoFresherXuanHoa")
		ctx.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func GoogleSignInCallBack(db *gorm.DB, oauth2cfg *oauth2.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sqlStorage := storage.NewSQLStore(db)
		userStorage := storage.NewUserStore(sqlStorage)
		authStorage := storage.NewAuthStore(userStorage)
		authBusiness := business.NewAuthBusiness(authStorage)

		// TODO: Random State to avoid CSRF attack:
		if state := ctx.Request.FormValue("state"); state != "dsoFresherXuanHoa" {
			fmt.Println("Error while get state from redirect url to authentication: state cannot be empty!")
			ctx.AbortWithStatus(http.StatusBadRequest)
		} else if token, err := oauth2cfg.Exchange(oauth2.NoContext, ctx.Request.FormValue("code")); err != nil {
			fmt.Println("Error while generate token to authentication: missing code or something wrong!")
			ctx.AbortWithStatus(http.StatusBadRequest)
		} else if res, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken); err != nil {
			fmt.Println("Error while require user information from Google authentication service: accessToken can be damaged or interrupt!")
			ctx.AbortWithStatus(http.StatusForbidden)
		} else {
			defer res.Body.Close()
			if usr, err := utils.NewOAuthUtil().OAuthResponse2User(res); err != nil {
				fmt.Println("Error while map OAuth response to User Creatable struct in auth transport: " + err.Error())
				ctx.JSON(http.StatusForbidden, entity.NewStandardResponse(nil, http.StatusForbidden, constants.StatusForbidden, err.Error(), constants.ErrMapOAuthResponse2UserCreatable))
			} else if accessToken, err := authBusiness.GoogleSignIn(ctx, usr); err != nil {
				fmt.Println("Error while sign in or sign up using Google account in auth transport: " + err.Error())
				ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), constants.ErrSignInUsingGoogle))
			} else {
				ctx.JSON(http.StatusOK, entity.NewStandardResponse(gin.H{"accessToken": accessToken}, http.StatusOK, constants.StatusOK, "", constants.SignInUsingGoogleSuccess))
			}

		}
	}
}

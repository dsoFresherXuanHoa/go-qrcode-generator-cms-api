package rest

import (
	"errors"
	"fmt"
	"go-qrcode-generator-cms-api/src/business"
	"go-qrcode-generator-cms-api/src/constants"
	"go-qrcode-generator-cms-api/src/entity"
	"go-qrcode-generator-cms-api/src/storage"
	"go-qrcode-generator-cms-api/src/utils"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

var (
	ErrMissingActivationCode = errors.New("missing activation code")
	ErrMissingEmail          = errors.New("missing email")

	InvalidSignUpRequestFormat            = "Invalid Sign Up Incoming Request: Check Swagger For More Information."
	ValidateSignUpRequestFailure          = "Invalid Sign Up Incoming Request: Check Swagger For More Information."
	EncodeImageFileMultiPartHeaderFailure = "Invalid Image File Type: Only Support PNG, JPG, JPEG."
	UploadSingleFileToCloudinaryFailure   = "Cannot Upload Single Image To Cloudinary: Only Support PNG, JPG, JPEG."
	SignUpFailure                         = "Cannot Sign Up An Account: Try Again Later."
	SignUpForNewUserSuccess               = "Sign Up An Account Success: Check Your Email To Activate Your Account."
	MissingActivationCode                 = "Missing Activation Code In Query String."
	ActivateAccountFailure                = "Cannot Activate Account: Try Again Later."
	ActivateAccountSuccess                = "Activate Account Success: Congrats."
	InvalidSignInRequestFormat            = "Invalid Sign In Incoming Request: Check Swagger For More Information."
	SignInFailure                         = "Cannot Sign In: Try Again Later."
	SignInSuccess                         = "Sign In Success: Congrats."
	GetDetailUserFailure                  = "Cannot Get Detail User: Try Again Later."
	GetDetailUserSuccess                  = "Get Detail User Information Success: Congrats."
	MissingEmail                          = "Missing Email In Query String."
	UserNotFound                          = "Cannot Get User Account With Email You Has Been Provide: Make Sure This Email Has Been Register Before."
	SendResetPasswordRequestEmailFailure  = "Cannot Send Reset Password Request Email: Congrats."
	SendResetPasswordRequestEmailSuccess  = "Send Reset Password Request Success: Check Your Email To Reset Your Password."
	InvalidResetPasswordRequestFormat     = "Invalid Reset Password Incoming Request: Check Swagger For More Information."
	ValidateResetPasswordRequestFailure   = "Invalid Reset Password Incoming Request: Check Swagger For More Information."
	ResetPasswordFailure                  = "Cannot Reset Password: Try Again Later."
	ResetPasswordSuccess                  = "Reset Password Success: Sign Out And Sign In Again!"
	ConvertOAuthResponse2UserFailure      = "Cannot Convert OAuth Response To User: Check Swagger For More Information."
	SignInUsingGoogleFailure              = "Cannot Sign In Using Google Account: Try Again Later."
	SignInUsingGoogleSuccess              = "Sign In Using Google Account: Congrats."
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
			fmt.Println("Error while parse user request: " + err.Error())
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, err.Error(), InvalidSignUpRequestFormat))
		} else if err := reqUser.Validate(); err != nil {
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, err.Error(), ValidateSignUpRequestFailure))
		} else if encodeAvatar, err := utils.NewImageUtil().ImageMultipartFile2Base64(reqUser.Avatar); err != nil {
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, err.Error(), EncodeImageFileMultiPartHeaderFailure))
		} else if uploadResult, err := cloudinaryStorage.UploadSingleImage(ctx, *encodeAvatar); err != nil {
			ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), UploadSingleFileToCloudinaryFailure))
		} else {
			reqUser.AvatarURL = uploadResult.URL
			if userUUID, err := authBusiness.SignUp(ctx, &reqUser); err != nil {
				fmt.Println("Error while sign up an account: " + err.Error())
				ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), SignUpFailure))
			} else {
				ctx.JSON(http.StatusOK, entity.NewStandardResponse(userUUID, http.StatusOK, "OK", "", SignUpForNewUserSuccess))
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
			fmt.Println("Error while get activationCode from user request: missing activation code")
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, ErrMissingActivationCode.Error(), MissingActivationCode))
		} else if err := authBusiness.Activate(ctx, activationCode); err != nil {
			fmt.Println("Error while activate user: " + err.Error())
			ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), ActivateAccountFailure))
		} else {
			ctx.JSON(http.StatusOK, entity.NewStandardResponse(true, http.StatusOK, constants.StatusOK, "", ActivateAccountSuccess))
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
			fmt.Println("Error while parse user request: " + err.Error())
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, err.Error(), InvalidSignInRequestFormat))
		} else if accessToken, err := authBusiness.SignIn(ctx, &reqUser); err != nil {
			fmt.Println("Error while sign in: " + err.Error())
			ctx.JSON(http.StatusUnauthorized, entity.NewStandardResponse(nil, http.StatusUnauthorized, constants.StatusUnauthorized, err.Error(), SignInFailure))
		} else {
			ctx.JSON(http.StatusOK, entity.NewStandardResponse(gin.H{"accessToken": accessToken}, http.StatusOK, constants.StatusOK, "", SignInSuccess))
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
			fmt.Println("Error while get detail user information: " + err.Error())
			ctx.JSON(http.StatusForbidden, entity.NewStandardResponse(nil, http.StatusForbidden, constants.StatusForbidden, err.Error(), GetDetailUserFailure))
		} else {
			ctx.JSON(http.StatusOK, entity.NewStandardResponse(user, http.StatusOK, constants.StatusOK, "", GetDetailUserSuccess))
		}
	}
}

// TODO: Do not direct use userStorage
func RequestResetPassword(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sqlStorage := storage.NewSQLStore(db)
		userStorage := storage.NewUserStore(sqlStorage)
		if email := ctx.Query("email"); email == "" {
			fmt.Println("Error while get email from user request: missing email")
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, ErrMissingEmail.Error(), MissingEmail))
		} else if usr, err := userStorage.FindUserByEmail(ctx, email); err != nil {
			ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), UserNotFound))
		} else if err := utils.NewMailUtil().SendResetPasswordRequestEmail(*usr.Email, usr.ActivationCode); err != nil {
			ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), SendResetPasswordRequestEmailFailure))
		} else {
			ctx.JSON(http.StatusOK, entity.NewStandardResponse(true, http.StatusOK, constants.StatusOK, "", SendResetPasswordRequestEmailSuccess))
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
			fmt.Println("Error while parse user request: " + err.Error())
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, err.Error(), InvalidResetPasswordRequestFormat))
		} else if err := reqUser.Validate(); err != nil {
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, err.Error(), ValidateResetPasswordRequestFailure))
		} else if resetCode := ctx.Query("resetCode"); resetCode == "" {
			fmt.Println("Error while get activationCode from user request: missing activation code")
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, ErrMissingActivationCode.Error(), MissingActivationCode))
		} else if err := authBusiness.ResetPassword(ctx, resetCode, &reqUser); err != nil {
			fmt.Println("Error while reset user password: " + err.Error())
			ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), ResetPasswordFailure))
		} else {
			ctx.JSON(http.StatusOK, entity.NewStandardResponse(true, http.StatusOK, constants.StatusOK, "", ResetPasswordSuccess))
		}
	}
}

func Home(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "sign-in.html", nil)
	}
}

// TODO: Random State to avoid CSRF attack:
func GoogleSignIn(db *gorm.DB, oauth2 *oauth2.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := oauth2.AuthCodeURL("dsoFresherXuanHoa")
		ctx.Redirect(http.StatusTemporaryRedirect, url)
	}
}

// TODO: Random State to avoid CSRF attack:
func GoogleSignInCallBack(db *gorm.DB, oauth2cfg *oauth2.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sqlStorage := storage.NewSQLStore(db)
		userStorage := storage.NewUserStore(sqlStorage)
		authStorage := storage.NewAuthStore(userStorage)
		authBusiness := business.NewAuthBusiness(authStorage)

		if state := ctx.Request.FormValue("state"); state != "dsoFresherXuanHoa" {
			fmt.Println("Error while get state from redirect url in order to authentication identity: state cannot be empty!")
			ctx.AbortWithStatus(http.StatusBadRequest)
		} else if token, err := oauth2cfg.Exchange(oauth2.NoContext, ctx.Request.FormValue("code")); err != nil {
			fmt.Println("Error while get code from redirect url in order to authentication identity: code cannot be empty!")
			ctx.AbortWithStatus(http.StatusBadRequest)
		} else if res, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken); err != nil {
			fmt.Println("Error while require user information from Google Authentication service: accessToken can be damaged or interrupt!")
			ctx.AbortWithStatus(http.StatusForbidden)
		} else {
			defer res.Body.Close()
			if usr, err := utils.NewOAuthUtil().Response2User(res); err != nil {
				fmt.Println("Error while get user info and convert http response: " + err.Error())
				ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), ConvertOAuthResponse2UserFailure))
			} else if accessToken, err := authBusiness.GoogleSignIn(ctx, usr); err != nil {
				fmt.Println("Error while sign in or sign up using Google account: " + err.Error())
				ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), SignInUsingGoogleFailure))
			} else {
				ctx.JSON(http.StatusOK, entity.NewStandardResponse(gin.H{"accessToken": accessToken}, http.StatusOK, constants.StatusOK, "", SignInUsingGoogleSuccess))
			}
		}
	}
}

package middlewares

import (
	"errors"
	"fmt"
	"go-qrcode-generator-cms-api/src/constants"
	"go-qrcode-generator-cms-api/src/entity"
	"go-qrcode-generator-cms-api/src/tokens/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ErrMissingToken           = errors.New("missing token from header")
	ErrMissingBearerToken     = errors.New("missing bearer token from header")
	ErrInvalidAccessToken     = errors.New("invalid access token")
	ErrPermissionDenied       = errors.New("permission denied")
	ErrPasswordHasBeenChanged = errors.New("your password has been changed")

	AdministratorPermission = 1
	UserPermission          = 2
)

func GetTokenFromHeader(c *gin.Context, key string) (token *string, err error) {
	if authHeader := c.Request.Header.Get(key); authHeader == "" {
		fmt.Println("Error while get Token from header: missing token ")
		return nil, ErrMissingToken
	} else if strings.Contains(authHeader, "Bearer ") {
		accessToken := strings.Split(authHeader, " ")[1]
		return &accessToken, nil
	} else {
		accessToken := authHeader
		return &accessToken, nil
	}
}

func RequiredAuthorized(db *gorm.DB, redisClient *redis.Client, secretKey string) gin.HandlerFunc {
	jwtTokenProvider := jwt.NewJWTProvider(secretKey)
	return func(c *gin.Context) {
		if authToken, err := GetTokenFromHeader(c, "Authorization"); err != nil {
			fmt.Println("Error while get Bearer token from header: " + err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, entity.NewStandardResponse(nil, http.StatusUnauthorized, constants.StatusUnauthorized, err.Error(), ErrMissingBearerToken.Error()))
		} else if jwtPayload, err := jwtTokenProvider.Validate(*authToken); err != nil {
			fmt.Println("Error while validate accessToken: " + err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, entity.NewStandardResponse(nil, http.StatusUnauthorized, constants.StatusUnauthorized, err.Error(), ErrInvalidAccessToken.Error()))
		} else if _, err := redisClient.Get(c, fmt.Sprint("accessTokenOfUser", jwtPayload.UserId)).Result(); err != nil {
			fmt.Println("Your password has been changed: " + err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, entity.NewStandardResponse(nil, http.StatusUnauthorized, constants.StatusUnauthorized, err.Error(), ErrPasswordHasBeenChanged.Error()))
		} else {
			userId := jwtPayload.UserId
			roleId := jwtPayload.RoleId
			c.Set("userId", userId)
			c.Set("roleId", roleId)
			c.Next()
		}
	}
}

func RequiredAdministratorPermission(db *gorm.DB, redisClient *redis.Client, secretKey string) gin.HandlerFunc {
	jwtTokenProvider := jwt.NewJWTProvider(secretKey)
	return func(c *gin.Context) {
		if authToken, err := GetTokenFromHeader(c, "Authorization"); err != nil {
			fmt.Println("Error while get Bearer token from header: " + err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, entity.NewStandardResponse(nil, http.StatusUnauthorized, constants.StatusUnauthorized, err.Error(), ErrMissingBearerToken.Error()))
		} else if jwtPayload, err := jwtTokenProvider.Validate(*authToken); err != nil {
			fmt.Println("Error while validate accessToken: " + err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, entity.NewStandardResponse(nil, http.StatusUnauthorized, constants.StatusUnauthorized, err.Error(), ErrInvalidAccessToken.Error()))
		} else if roleId := jwtPayload.RoleId; int(roleId) != AdministratorPermission {
			fmt.Println("Error while validate user permission: you don't has right permission to do this.")
			c.AbortWithStatusJSON(http.StatusForbidden, entity.NewStandardResponse(nil, http.StatusForbidden, constants.StatusForbidden, ErrPermissionDenied.Error(), ErrPermissionDenied.Error()))
		} else if _, err := redisClient.Get(c, fmt.Sprint("accessTokenOfUser", jwtPayload.UserId)).Result(); err != nil {
			fmt.Println("Your password has been changed: " + err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, entity.NewStandardResponse(nil, http.StatusUnauthorized, constants.StatusUnauthorized, err.Error(), ErrPasswordHasBeenChanged.Error()))
		} else {
			userId := jwtPayload.UserId
			roleId := jwtPayload.RoleId
			c.Set("userId", userId)
			c.Set("roleId", roleId)
			c.Next()
		}
	}
}

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
	"gorm.io/gorm"
)

var (
	ErrMissingToken       = errors.New("missing token from header")
	ErrMissingBearerToken = errors.New("missing bearer token from header")
	ErrInvalidAccessToken = errors.New("invalid access token")
)

func GetTokenFromHeader(c *gin.Context, key string) (token *string, err error) {
	if authHeader := c.Request.Header.Get(key); authHeader == "" {
		fmt.Println("Error while get Token from header: missing token ")
		return nil, ErrMissingToken
	} else {
		accessToken := strings.Split(authHeader, " ")[1]
		return &accessToken, nil
	}
}

func RequiredAuthorized(db *gorm.DB, secretKey string) gin.HandlerFunc {
	jwtTokenProvider := jwt.NewJWTProvider(secretKey)
	return func(c *gin.Context) {
		if authToken, err := GetTokenFromHeader(c, "Authorization"); err != nil {
			fmt.Println("Error while get Bearer token from header: " + err.Error())
			c.JSON(http.StatusUnauthorized, entity.NewStandardResponse(nil, http.StatusUnauthorized, constants.StatusUnauthorized, err.Error(), ErrMissingBearerToken.Error()))
		} else if jwtPayload, err := jwtTokenProvider.Validate(*authToken); err != nil {
			fmt.Println("Error while validate accessToken: " + err.Error())
			c.JSON(http.StatusUnauthorized, entity.NewStandardResponse(nil, http.StatusUnauthorized, constants.StatusUnauthorized, err.Error(), ErrInvalidAccessToken.Error()))
		} else {
			userId := jwtPayload.UserId
			roleId := jwtPayload.RoleId
			c.Set("userId", userId)
			c.Set("roleId", roleId)
			c.Next()
		}
	}
}

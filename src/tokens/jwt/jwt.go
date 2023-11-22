package jwt

import (
	"errors"
	"fmt"
	"go-qrcode-generator-cms-api/src/tokens"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrGenerateAccessToken       = errors.New("generate accessToken failure")
	ErrValidateAccessToken       = errors.New("validate accessToken failure")
	ErrClaimsAccessToken2Payload = errors.New("claims accessToken to payload failure")
)

type jwtProvider struct {
	secretKey string
}

func NewJWTProvider(secretKey string) *jwtProvider {
	return &jwtProvider{secretKey: secretKey}
}

type authClaims struct {
	Payload tokens.TokenPayload `json:"tokenPayload"`
	jwt.StandardClaims
}

func (j *jwtProvider) Generate(payload tokens.TokenPayload, exp int) (*tokens.Token, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims{
		Payload: payload,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(exp)).Unix(),
			IssuedAt:  time.Now().Local().Unix(),
		},
	})
	if accessToken, err := t.SignedString([]byte(j.secretKey)); err != nil {
		fmt.Println("Error while generate accessToken: " + err.Error())
		return nil, ErrGenerateAccessToken
	} else {
		return &tokens.Token{Token: accessToken, AvailableUntil: exp, CreatedAt: time.Now()}, nil
	}
}

func (j *jwtProvider) Validate(token string) (*tokens.TokenPayload, error) {
	if res, err := jwt.ParseWithClaims(token, &authClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	}); err != nil {
		fmt.Println("Error while validate accessToken: " + err.Error())
		return nil, ErrValidateAccessToken
	} else if claims, ok := res.Claims.(*authClaims); !ok {
		fmt.Println("Error while claims accessToken to payload: " + err.Error())
		return nil, ErrClaimsAccessToken2Payload
	} else {
		return &claims.Payload, nil
	}
}

package jwt

import (
	"fmt"
	"go-qrcode-generator-cms-api/src/tokens"
	"time"

	"github.com/golang-jwt/jwt"
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
			ExpiresAt: time.Now().Local().Add(time.Second * time.Duration(exp)).Unix(),
			IssuedAt:  time.Now().Local().Unix(),
		},
	})
	if accessToken, err := t.SignedString([]byte(j.secretKey)); err != nil {
		fmt.Println("Error while create accessToken in jwt provider: " + err.Error())
		return nil, err
	} else {
		return &tokens.Token{Token: accessToken, AvailableUntil: exp, CreatedAt: time.Now()}, nil
	}
}

func (j *jwtProvider) Validate(token string) (*tokens.TokenPayload, error) {
	if res, err := jwt.ParseWithClaims(token, &authClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	}); err != nil {
		fmt.Println("Error while validate accessToken in jwt provider: " + err.Error())
		return nil, err
	} else if claims, ok := res.Claims.(*authClaims); !ok {
		fmt.Println("Error while validate authClaims struct in jwt provider: " + err.Error())
		return nil, err
	} else {
		return &claims.Payload, nil
	}
}

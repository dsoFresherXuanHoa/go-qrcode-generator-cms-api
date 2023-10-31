package business

import (
	"context"
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"
	"go-qrcode-generator-cms-api/src/tokens"
	"go-qrcode-generator-cms-api/src/tokens/jwt"
	"os"
)

type AuthStorage interface {
	SignUp(ctx context.Context, user *entity.UserCreatable) (*string, error)
	Activate(ctx context.Context, activationCode string) error
	SignIn(ctx context.Context, user *entity.UserQueryable) (*entity.UserResponse, error)
	Me(ctx context.Context, userId uint) (*entity.UserResponse, error)
}

type authBusiness struct {
	authStorage AuthStorage
}

func NewAuthBusiness(authStorage AuthStorage) *authBusiness {
	return &authBusiness{authStorage: authStorage}
}

func (business *authBusiness) SignUp(ctx context.Context, user *entity.UserCreatable) (*string, error) {
	if err := user.Validate(); err != nil {
		fmt.Println("Error while validate user request in sign up business: " + err.Error())
		return nil, err
	} else {
		user.Mask()
		if userUUID, err := business.authStorage.SignUp(ctx, user); err != nil {
			fmt.Println("Error while save user information to database in sign up business: " + err.Error())
			return nil, err
		} else {
			return userUUID, nil
		}
	}
}

func (business *authBusiness) Activate(ctx context.Context, activationCode string) error {
	if err := business.authStorage.Activate(ctx, activationCode); err != nil {
		fmt.Println("Error while active account by activation code in auth business: " + err.Error())
		return err
	}
	return nil
}

func (business *authBusiness) SignIn(ctx context.Context, user *entity.UserQueryable) (*tokens.Token, error) {
	if err := user.Validate(); err != nil {
		fmt.Println("Error while validate user request in sign sin business: " + err.Error())
		return nil, err
	} else if usr, err := business.authStorage.SignIn(ctx, user); err != nil {
		fmt.Println("Error while find user by email and password: " + err.Error())
		return nil, err
	} else {
		secretKey := os.Getenv("JWT_ACCESS_SECRET")
		payload := tokens.TokenPayload{UserId: usr.ID, RoleId: usr.Role.ID}
		jwtProvider := jwt.NewJWTProvider(secretKey)
		if accessToken, err := jwtProvider.Generate(payload, 86400); err != nil {
			fmt.Println("Error while try to generate accessToken in auth business: " + err.Error())
			return nil, err
		} else {
			return accessToken, nil
		}
	}
}

func (business *authBusiness) Me(ctx context.Context, userId uint) (*entity.UserResponse, error) {
	if usr, err := business.authStorage.Me(ctx, userId); err != nil {
		fmt.Println("Error while find detail user by id (hidden id) in auth business: " + err.Error())
		return nil, err
	} else {
		return usr, nil
	}
}

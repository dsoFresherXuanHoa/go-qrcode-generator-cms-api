package business

import (
	"context"
	"go-qrcode-generator-cms-api/src/entity"
	"go-qrcode-generator-cms-api/src/tokens"
	"go-qrcode-generator-cms-api/src/tokens/jwt"
	"go-qrcode-generator-cms-api/src/utils"
	"os"
)

type AuthStorage interface {
	SignUp(ctx context.Context, user *entity.UserCreatable) (*string, error)
	Activate(ctx context.Context, activationCode string) error
	SignIn(ctx context.Context, user *entity.UserQueryable) (*entity.UserResponse, error)
	Me(ctx context.Context, userId uint) (*entity.UserResponse, error)
	ResetPassword(ctx context.Context, activationCode string, user *entity.UserUpdatable) error
	GoogleSignUp(ctx context.Context, user *entity.UserCreatable) (*string, error)
	GoogleSignIn(ctx context.Context, email string) (*entity.UserResponse, error)
	VerifyEmailHasBeenUsed(ctx context.Context, email string) (bool, error)
}

type authBusiness struct {
	authStorage AuthStorage
}

func NewAuthBusiness(authStorage AuthStorage) *authBusiness {
	return &authBusiness{authStorage: authStorage}
}

func (business *authBusiness) SignUp(ctx context.Context, user *entity.UserCreatable) (*string, error) {
	user.Mask()
	if userUUID, err := business.authStorage.SignUp(ctx, user); err != nil {
		return nil, err
	} else if err := utils.NewMailUtil().SendActivationRequestEmail(*user.Email, user.ActivationCode); err != nil {
		return nil, err
	} else {
		return userUUID, nil
	}
}

func (business *authBusiness) Activate(ctx context.Context, activationCode string) error {
	if err := business.authStorage.Activate(ctx, activationCode); err != nil {
		return err
	}
	return nil
}

func (business *authBusiness) SignIn(ctx context.Context, user *entity.UserQueryable) (*tokens.Token, error) {
	if usr, err := business.authStorage.SignIn(ctx, user); err != nil {
		return nil, err
	} else {
		secretKey := os.Getenv("JWT_ACCESS_SECRET")
		payload := tokens.TokenPayload{UserId: usr.ID, RoleId: usr.Role.ID}
		jwtProvider := jwt.NewJWTProvider(secretKey)
		if accessToken, err := jwtProvider.Generate(payload, 86400); err != nil {
			return nil, err
		} else {
			return accessToken, nil
		}
	}
}

func (business *authBusiness) Me(ctx context.Context, userId uint) (*entity.UserResponse, error) {
	if usr, err := business.authStorage.Me(ctx, userId); err != nil {
		return nil, err
	} else {
		return usr, nil
	}
}

func (business *authBusiness) ResetPassword(ctx context.Context, activationCode string, user *entity.UserUpdatable) error {
	user.Mask()
	if err := business.authStorage.ResetPassword(ctx, activationCode, user); err != nil {
		return err
	}
	return nil
}

func (business *authBusiness) GoogleSignIn(ctx context.Context, user *entity.UserCreatable) (*tokens.Token, error) {
	if _, err := business.authStorage.VerifyEmailHasBeenUsed(ctx, *user.Email); err != nil {
		if _, err := business.authStorage.GoogleSignUp(ctx, user); err != nil {
			return nil, err
		}
	}

	if usr, err := business.authStorage.GoogleSignIn(ctx, *user.Email); err != nil {
		return nil, err
	} else {
		secretKey := os.Getenv("JWT_ACCESS_SECRET")
		payload := tokens.TokenPayload{UserId: usr.ID, RoleId: usr.Role.ID}
		jwtProvider := jwt.NewJWTProvider(secretKey)
		if accessToken, err := jwtProvider.Generate(payload, 86400); err != nil {
			return nil, err
		} else {
			return accessToken, nil
		}
	}
}

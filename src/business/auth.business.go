package business

import (
	"context"
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"
)

type AuthStorage interface {
	SignUp(ctx context.Context, user *entity.UserCreatable) (*string, error)
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

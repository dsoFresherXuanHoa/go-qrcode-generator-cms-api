package storage

import (
	"context"
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"
)

type authStorage struct {
	userStorage *userStorage
}

func NewAuthStore(userStorage *userStorage) *authStorage {
	return &authStorage{userStorage: userStorage}
}

func (s *authStorage) SignUp(ctx context.Context, user *entity.UserCreatable) (*string, error) {
	if userUUID, err := s.userStorage.CreateUser(ctx, user); err != nil {
		fmt.Println("Error while save user information to database in auth storage: " + err.Error())
		return nil, err
	} else {
		return userUUID, nil
	}
}

func (s *authStorage) Activate(ctx context.Context, activationCode string) error {
	if err := s.userStorage.UpdateUserActivateStatusByActivationCode(ctx, activationCode); err != nil {
		fmt.Println("Error while activate user by activation code in auth storage: " + err.Error())
		return err
	}
	return nil
}

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

func (s *authStorage) SignIn(ctx context.Context, user *entity.UserQueryable) (*entity.UserResponse, error) {
	if usr, err := s.userStorage.FindUserByEmailAndPassword(ctx, *user.Email, *user.Password); err != nil {
		fmt.Println("Error while find user by email and password in auth storage: " + err.Error())
		return nil, err
	} else {
		return usr, nil
	}
}

func (s *authStorage) Me(ctx context.Context, userId uint) (*entity.UserResponse, error) {
	if usr, err := s.userStorage.FindDetailUserById(ctx, userId); err != nil {
		fmt.Println("Error while find detail user by id (hidden id) in auth storage: " + err.Error())
		return nil, err
	} else {
		return usr, nil
	}
}

func (s *authStorage) ResetPassword(ctx context.Context, activationCode string, user *entity.UserUpdatable) error {
	if err := s.userStorage.UpdateUserPasswordByActivationCode(ctx, activationCode, user); err != nil {
		fmt.Println("Error while update user password in auth storage: " + err.Error())
	}
	return nil
}

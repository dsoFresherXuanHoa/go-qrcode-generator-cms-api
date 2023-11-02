package storage

import (
	"context"
	"errors"
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"
)

var (
	ErrSignUp                       = errors.New("sign up for an user failure")
	ErrActivateUserByActivationCode = errors.New("activate user by activation code failure")
	ErrSignInUsingEmailPassword     = errors.New("sign in using email and password failure")
	ErrGetUsingInfo                 = errors.New("get user information failure")
	ErrResetPassword                = errors.New("reset user password by activation code failure")
	ErrGoogleSignUp                 = errors.New("sign up for new user using Google account failure")
	ErrGoogleSignIn                 = errors.New("sign in using Google account failure")
	ErrEmailHasBeenUsed             = errors.New("verify email failure: this email has been used")
)

type authStorage struct {
	userStorage *userStorage
}

func NewAuthStore(userStorage *userStorage) *authStorage {
	return &authStorage{userStorage: userStorage}
}

func (s *authStorage) SignUp(ctx context.Context, user *entity.UserCreatable) (*string, error) {
	if userUUID, err := s.userStorage.CreateUser(ctx, user); err != nil {
		fmt.Println("Error while sign up for an user: " + err.Error())
		return nil, ErrSignUp
	} else {
		return userUUID, nil
	}
}

func (s *authStorage) Activate(ctx context.Context, activationCode string) error {
	if err := s.userStorage.UpdateUserActivateStatusByActivationCode(ctx, activationCode); err != nil {
		fmt.Println("Error while activate user by activation code: " + err.Error())
		return ErrActivateUserByActivationCode
	}
	return nil
}

func (s *authStorage) SignIn(ctx context.Context, user *entity.UserQueryable) (*entity.UserResponse, error) {
	if usr, err := s.userStorage.FindUserByEmailAndPassword(ctx, *user.Email, *user.Password); err != nil {
		fmt.Println("Error while sign in using email and password: " + err.Error())
		return nil, ErrSignInUsingEmailPassword
	} else {
		return usr, nil
	}
}

func (s *authStorage) Me(ctx context.Context, userId uint) (*entity.UserResponse, error) {
	if usr, err := s.userStorage.FindDetailUserById(ctx, userId); err != nil {
		fmt.Println("Error while get user information: " + err.Error())
		return nil, ErrGetUsingInfo
	} else {
		return usr, nil
	}
}

func (s *authStorage) ResetPassword(ctx context.Context, activationCode string, user *entity.UserUpdatable) error {
	if err := s.userStorage.UpdateUserPasswordByActivationCode(ctx, activationCode, user); err != nil {
		fmt.Println("Error while reset user password by activation code: " + err.Error())
		return ErrResetPassword
	}
	return nil
}

// TODO: Do not mask and set activate status in storage
func (s *authStorage) GoogleSignUp(ctx context.Context, user *entity.UserCreatable) (*string, error) {
	user.Activate = true
	user.Mask()
	if userUUID, err := s.userStorage.CreateUser(ctx, user); err != nil {
		fmt.Println("Error while sign up for new user using Google account: " + err.Error())
		return nil, ErrGoogleSignUp
	} else {
		return userUUID, nil
	}
}

func (s *authStorage) GoogleSignIn(ctx context.Context, email string) (*entity.UserResponse, error) {
	if usr, err := s.userStorage.FindUserByEmail(ctx, email); err != nil {
		fmt.Println("Error while sign in using Google account: " + err.Error())
		return nil, ErrGoogleSignIn
	} else {
		res := usr.Convert2Response()
		return &res, nil
	}
}

func (s *authStorage) VerifyEmailHasBeenUsed(ctx context.Context, email string) (bool, error) {
	if _, err := s.userStorage.FindUserByEmail(ctx, email); err != nil {
		fmt.Println("Error while verify email had not been used: " + err.Error())
		return false, ErrEmailHasBeenUsed
	}
	return true, nil
}

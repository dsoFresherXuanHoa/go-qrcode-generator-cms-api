package storage

import (
	"context"
	"errors"
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrSaveUser                                 = errors.New("save user failure")
	ErrUpdateUserActivateStatusByActivationCode = errors.New("update user activate status by activation code failure")
	ErrFindUserByEmailAndActivationCode         = errors.New("find user by email and activate status failure")
	ErrCompareUserPassword                      = errors.New("compare user password and hash password failure")
	ErrFindUserById                             = errors.New("find detail user by id failure")
	ErrUpdateUserPasswordByActivationCode       = errors.New("update user password by activation code failure")
	ErrFindUserByEmail                          = errors.New("find user by email failure")
)

type userStorage struct {
	sql *sqlStorage
}

func NewUserStore(sql *sqlStorage) *userStorage {
	return &userStorage{sql: sql}
}

func (s *userStorage) CreateUser(ctx context.Context, user *entity.UserCreatable) (*string, error) {
	if err := s.sql.db.Table(entity.UserCreatable{}.TableName()).Create(&user).Error; err != nil {
		fmt.Println("Error while save user into database: " + err.Error())
		return nil, ErrSaveRole2DB
	}
	return &user.UUID, nil
}

func (s *userStorage) UpdateUserActivateStatusByActivationCode(ctx context.Context, activationCode string) error {
	if err := s.sql.db.Model(&entity.User{}).Where("activation_code = ?", activationCode).Update("activate", true).Error; err != nil {
		fmt.Println("Error while update user activate status by activation code into database: " + err.Error())
		return ErrUpdateUserActivateStatusByActivationCode
	}
	return nil
}

func (s *userStorage) FindUserByEmailAndPassword(ctx context.Context, email string, password string) (*entity.UserResponse, error) {
	var usr entity.UserResponse
	if err := s.sql.db.Table(usr.TableName()).Where("email = ?", email).Where("activate = ?", true).First(&usr).Error; err != nil {
		fmt.Println("Error while get user by email and activation status from database: " + err.Error())
		return nil, ErrFindUserByEmailAndActivationCode
	} else if err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password)); err != nil {
		fmt.Println("Error while compare user password and hash password from database: " + err.Error())
		return nil, ErrCompareUserPassword
	}
	return &usr, nil
}

func (s *userStorage) FindDetailUserById(ctx context.Context, id uint) (*entity.UserResponse, error) {
	var usr entity.UserResponse
	if err := s.sql.db.Table(usr.TableName()).Where("id = ?", id).Preload("Role").First(&usr).Error; err != nil {
		fmt.Println("Error while find detail user by id from database: " + err.Error())
		return nil, ErrFindUserById
	}
	return &usr, nil
}

func (s *userStorage) UpdateUserPasswordByActivationCode(ctx context.Context, activationCode string, user *entity.UserUpdatable) error {
	if err := s.sql.db.Model(&entity.User{}).Where("activation_code = ?", activationCode).Update("password", *user.Password).Error; err != nil {
		fmt.Println("Error while update user password by activation code into database: " + err.Error())
		return ErrUpdateUserPasswordByActivationCode
	}
	return nil
}

func (s *userStorage) FindUserByEmail(ctx context.Context, email string) (*entity.UserCreatable, error) {
	var usr entity.UserCreatable
	if err := s.sql.db.Table(usr.TableName()).Where("email = ?", email).Where("activate = ?", true).First(&usr).Error; err != nil {
		fmt.Println("Error while find user by email from database: " + err.Error())
		return nil, ErrFindUserByEmail
	}
	return &usr, nil
}

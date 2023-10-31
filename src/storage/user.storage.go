package storage

import (
	"context"
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"

	"golang.org/x/crypto/bcrypt"
)

type userStorage struct {
	sql *sqlStorage
}

func NewUserStore(sql *sqlStorage) *userStorage {
	return &userStorage{sql: sql}
}

func (s *userStorage) CreateUser(ctx context.Context, user *entity.UserCreatable) (*string, error) {
	if err := s.sql.db.Table(entity.UserCreatable{}.GetTableName()).Create(&user).Error; err != nil {
		fmt.Println("Error while save user information to database in user storage: " + err.Error())
		return nil, err
	}
	return &user.UUID, nil
}

func (s *userStorage) UpdateUserActivateStatusByActivationCode(ctx context.Context, activationCode string) error {
	if err := s.sql.db.Model(&entity.User{}).Where("activation_code = ?", activationCode).Update("activate", true).Error; err != nil {
		fmt.Println("Error while activate user by activation code in user storage: " + err.Error())
		return err
	}
	return nil
}

func (s *userStorage) FindUserByEmailAndPassword(ctx context.Context, email string, password string) (*entity.UserResponse, error) {
	var usr entity.UserResponse
	if err := s.sql.db.Table(usr.GetTableName()).Where("email = ?", email).Where("activate = ?", true).First(&usr).Error; err != nil {
		fmt.Println("Error while find user by email and password in user storage: " + err.Error())
		return nil, err
	} else if err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password)); err != nil {
		fmt.Println("Error while compare user password and hash password in user storage: " + err.Error())
		return nil, err
	}
	return &usr, nil
}

func (s *userStorage) FindDetailUserById(ctx context.Context, id uint) (*entity.UserResponse, error) {
	var usr entity.UserResponse
	if err := s.sql.db.Table(usr.GetTableName()).Where("id = ?", id).Preload("Role").First(&usr).Error; err != nil {
		fmt.Println("Error while find detail user by id (hidden id) in user storage: " + err.Error())
		return nil, err
	}
	return &usr, nil
}

func (s *userStorage) UpdateUserPasswordByActivationCode(ctx context.Context, activationCode string, user *entity.UserUpdatable) error {
	if err := s.sql.db.Model(&entity.User{}).Where("activation_code = ?", activationCode).Update("password", *user.Password).Error; err != nil {
		fmt.Println("Error while reset user password by activation code in user storage: " + err.Error())
		return err
	}
	return nil
}

func (s *userStorage) FindUserByEmail(ctx context.Context, email string) (*entity.UserCreatable, error) {
	var usr entity.UserCreatable
	if err := s.sql.db.Table(usr.GetTableName()).Where("email = ?", email).Where("activate = ?", true).First(&usr).Error; err != nil {
		fmt.Println("Error while find user by email in user storage: " + err.Error())
		return nil, err
	}
	return &usr, nil
}

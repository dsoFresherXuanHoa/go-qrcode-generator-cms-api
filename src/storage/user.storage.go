package storage

import (
	"context"
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"
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

package storage

import (
	"context"
	"errors"
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"
)

var (
	ErrSaveRole2DB = errors.New("save role into database failure")
)

type roleStorage struct {
	sql *sqlStorage
}

func NewRoleStore(sql *sqlStorage) *roleStorage {
	return &roleStorage{sql: sql}
}

func (s *roleStorage) CreateRole(ctx context.Context, role *entity.RoleCreatable) (*string, error) {
	if err := s.sql.db.Table(entity.RoleCreatable{}.TableName()).Create(&role).Error; err != nil {
		fmt.Println("Error while save role into database: " + err.Error())
		return nil, ErrSaveRole2DB
	}
	return &role.UUID, nil
}

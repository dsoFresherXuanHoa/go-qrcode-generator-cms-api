package storage

import (
	"context"
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"
)

type roleStorage struct {
	sql *sqlStorage
}

func NewRoleStore(sql *sqlStorage) *roleStorage {
	return &roleStorage{sql: sql}
}

func (s *roleStorage) CreateRole(ctx context.Context, role *entity.RoleCreatable) (*string, error) {
	if err := s.sql.db.Table(entity.RoleCreatable{}.GetTableName()).Create(&role).Error; err != nil {
		fmt.Println("Error while save role information to database in role storage: " + err.Error())
		return nil, err
	}
	return &role.UUID, nil
}

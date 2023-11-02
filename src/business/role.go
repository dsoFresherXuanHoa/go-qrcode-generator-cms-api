package business

import (
	"context"
	"go-qrcode-generator-cms-api/src/entity"
)

type RoleStorage interface {
	CreateRole(ctx context.Context, role *entity.RoleCreatable) (*string, error)
}

type roleBusiness struct {
	roleStorage RoleStorage
}

func NewRoleBusiness(roleStorage RoleStorage) *roleBusiness {
	return &roleBusiness{roleStorage: roleStorage}
}

func (business *roleBusiness) CreateRole(ctx context.Context, role *entity.RoleCreatable) (*string, error) {
	role.Mask()
	if roleUUID, err := business.roleStorage.CreateRole(ctx, role); err != nil {
		return nil, err
	} else {
		return roleUUID, nil
	}
}

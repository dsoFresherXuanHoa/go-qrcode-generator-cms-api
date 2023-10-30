package business

import (
	"context"
	"fmt"
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
	if err := role.Validate(); err != nil {
		fmt.Println("Error while validate user request in role business: " + err.Error())
		return nil, err
	} else {
		role.Mask()
		if roleUUID, err := business.roleStorage.CreateRole(ctx, role); err != nil {
			fmt.Println("Error while save role information to database in role business: " + err.Error())
			return nil, err
		} else {
			return roleUUID, nil
		}
	}
}

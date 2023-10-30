package entity

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model `json:"-"`
	User       Users  `json:"-"`
	UUID       string `json:"-" gorm:"not null"`
	Name       string `json:"-" gorm:"not null"`
}

type RoleResponse struct {
	gorm.Model `json:"-"`
	UUID       string `json:"uuid" gorm:"not null"`
	Name       string `json:"name" gorm:"not null"`
}

// TODO: Custom error message
type RoleCreatable struct {
	gorm.Model `json:"-"`
	UUID       string  `json:"-" gorm:"not null"`
	Name       *string `json:"name" validate:"required" gorm:"not null"`
}

type RoleUpdatable struct {
	gorm.Model `json:"-"`
	Name       *string `json:"name" validate:"required" gorm:"not null"`
}

type Roles []Role

func (Role) GetTableName() string          { return "roles" }
func (Roles) GetTableName() string         { return Role{}.GetTableName() }
func (RoleCreatable) GetTableName() string { return Role{}.GetTableName() }
func (RoleUpdatable) GetTableName() string { return Role{}.GetTableName() }

func (role RoleCreatable) Validate() error {
	validate := validator.New()
	if err := validate.Struct(&role); err != nil {
		fmt.Println("Error while validate role creatable: " + err.Error())
		return err
	}
	return nil
}

func (role RoleUpdatable) Validate() error {
	validate := validator.New()
	if err := validate.Struct(&role); err != nil {
		fmt.Println("Error while validate user updatable: " + err.Error())
		return err
	}
	return nil
}

func (role *RoleCreatable) Mask() {
	role.UUID = uuid.NewString()
}

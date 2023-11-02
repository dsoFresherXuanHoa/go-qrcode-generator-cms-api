package entity

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	User Users

	UUID string `gorm:"not null"`
	Name string `gorm:"not null"`
}

type RoleResponse struct {
	gorm.Model `json:"-"`
	Users      Users `json:"users" gorm:"foreignKey:RoleId"`

	UUID string `json:"uuid" gorm:"not null"`
	Name string `json:"name" gorm:"not null"`
}

type RoleCreatable struct {
	gorm.Model

	UUID string  `json:"-" gorm:"not null"`
	Name *string `json:"name" validate:"required" gorm:"not null"`
}

type RoleUpdatable struct {
	gorm.Model `json:"-"`

	Name *string `json:"name" validate:"required" gorm:"not null"`
}

type Roles []Role

func (Role) TableName() string          { return "roles" }
func (Roles) TableName() string         { return Role{}.TableName() }
func (RoleCreatable) TableName() string { return Role{}.TableName() }
func (RoleUpdatable) TableName() string { return Role{}.TableName() }

func (role RoleCreatable) Validate() error {
	validate := validator.New()
	if err := validate.Struct(&role); err != nil {
		fmt.Println("Error while validate incoming request role: " + err.Error())
		return err
	}
	return nil
}

func (role RoleUpdatable) Validate() error {
	validate := validator.New()
	if err := validate.Struct(&role); err != nil {
		fmt.Println("Error while validate incoming request role: " + err.Error())
		return err
	}
	return nil
}

func (role *RoleCreatable) Mask() {
	role.UUID = uuid.NewString()
}

package entity

import (
	"fmt"
	"mime/multipart"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model     `json:"-"`
	RoleID         uint   `json:"-"`
	UUID           string `json:"-" gorm:"not null"`
	FirstName      string `json:"-" gorm:"column:first_name;not null"`
	LastName       string `json:"-" gorm:"column:last_name;not null"`
	Gender         bool   `json:"-" gorm:"default:false"`
	Email          string `json:"-" gorm:"unique;not null"`
	Password       string `json:"-" gorm:"not null"`
	Activate       bool   `json:"-" gorm:"default:false"`
	ActivationCode string `json:"-" gorm:"column:activation_code;not null"`
	AvatarURL      string `json:"-" gorm:"column:avatar_url;not null"`
}

type UserResponse struct {
	gorm.Model     `json:"-"`
	UUID           string `json:"uuid" gorm:"not null"`
	FirstName      string `json:"firstName" gorm:"column:first_name;not null"`
	LastName       string `json:"lastName" gorm:"column:last_name;not null"`
	Gender         bool   `json:"gender" gorm:"default:false"`
	Email          string `json:"email" gorm:"unique;not null"`
	Password       string `json:"-" gorm:"not null"`
	Activate       bool   `json:"activate" gorm:"default:false"`
	ActivationCode string `json:"-" gorm:"column:activation_code;not null"`
	AvatarURL      string `json:"avatarUrl" gorm:"column:avatar_url;not null"`
}

// TODO: Custom error message
type UserCreatable struct {
	gorm.Model     `json:"-"`
	RoleID         *uint                 `form:"roleId" json:"roleId" validate:"required" gorm:"column:role_id;not null"`
	UUID           string                `form:"-" json:"-" gorm:"not null"`
	FirstName      *string               `form:"firstName" json:"firstName" validate:"required" gorm:"column:first_name;not null"`
	LastName       *string               `form:"lastName" json:"lastName" validate:"required" gorm:"column:last_name;not null"`
	Gender         *bool                 `form:"gender" json:"gender" validate:"required" gorm:"default:false"`
	Email          *string               `form:"email" json:"email" validate:"required,email" gorm:"unique;not null"`
	Password       *string               `form:"password" json:"password" validate:"required,min=8" gorm:"not null"`
	Activate       bool                  `form:"-" json:"-" gorm:"default:false"`
	ActivationCode string                `form:"-" json:"-" gorm:"activation_code;not null"`
	Avatar         *multipart.FileHeader `form:"avatar" validate:"required" sql:"-" gorm:"-"`
	AvatarURL      string                `form:"-" json:"-" gorm:"column:avatar_url;not null"`
}

type UserUpdatable struct {
	gorm.Model `json:"-"`
	Password   *string `json:"password" validate:"required,min=8" gorm:"not null"`
	Activate   bool    `json:"-" gorm:"default:false"`
}

type Users []User

func (User) GetTableName() string          { return "users" }
func (Users) GetTableName() string         { return User{}.GetTableName() }
func (UserCreatable) GetTableName() string { return User{}.GetTableName() }
func (UserUpdatable) GetTableName() string { return User{}.GetTableName() }

func (usr UserCreatable) Validate() error {
	validate := validator.New()
	if err := validate.Struct(&usr); err != nil {
		fmt.Println("Error while validate user creatable: " + err.Error())
		return err
	}
	return nil
}

func (usr UserUpdatable) Validate() error {
	validate := validator.New()
	if err := validate.Struct(&usr); err != nil {
		fmt.Println("Error while validate user creatable: " + err.Error())
		return err
	}
	return nil
}

func (usr *UserCreatable) Mask() {
	hashPasswordBytes, _ := bcrypt.GenerateFromPassword([]byte(*usr.Password), 5)
	hashPassword := string(hashPasswordBytes)

	*usr.Password = hashPassword
	usr.UUID = uuid.NewString()
	usr.ActivationCode = uuid.NewString()
}

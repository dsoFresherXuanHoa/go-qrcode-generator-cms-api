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
	gorm.Model
	RoleID  uint
	QrCodes QRCodes

	UUID           string `gorm:"not null"`
	FirstName      string `gorm:"not null"`
	LastName       string `gorm:"not null"`
	Gender         bool   `gorm:"default:false"`
	Email          string `gorm:"unique;not null"`
	Password       string `gorm:"not null"`
	Activate       bool   `gorm:"default:false"`
	ActivationCode string `gorm:"not null"`
	AvatarURL      string `gorm:"not null"`
}

type UserResponse struct {
	gorm.Model `json:"-"`
	RoleId     uint `json:"-"`
	Role       Role `json:"role"`

	UUID           string `json:"uuid" gorm:"not null"`
	FirstName      string `json:"firstName" gorm:"not null"`
	LastName       string `json:"lastName" gorm:"not null"`
	Gender         bool   `json:"gender" gorm:"default:false"`
	Email          string `json:"email" gorm:"unique;not null"`
	Activate       bool   `json:"activate" gorm:"default:false"`
	AvatarURL      string `json:"avatarUrl" gorm:"not null"`
	ActivationCode string `json:"-" gorm:"not null"`
	Password       string `json:"-" gorm:"not null"`
}

type UserCreatable struct {
	gorm.Model

	RoleID         *uint   `form:"roleId" json:"roleId" validate:"required" gorm:"not null"`
	FirstName      *string `form:"firstName" json:"firstName" validate:"required" gorm:"not null"`
	LastName       *string `form:"lastName" json:"lastName" validate:"required" gorm:"not null"`
	Gender         *bool   `form:"gender" json:"gender" validate:"required" gorm:"default:false"`
	Email          *string `form:"email" json:"email" validate:"required,email" gorm:"unique;not null"`
	Password       *string `form:"password" json:"password" validate:"required,min=8" gorm:"not null"`
	UUID           string  `form:"-" json:"-" gorm:"not null"`
	Activate       bool    `form:"-" json:"-" gorm:"default:false"`
	ActivationCode string  `form:"-" json:"-" gorm:"not null"`
	AvatarURL      string  `form:"-" json:"-" gorm:"not null"`

	Avatar *multipart.FileHeader `form:"avatar" validate:"required" sql:"-" gorm:"-"`
}

type UserUpdatable struct {
	gorm.Model
	Password *string `json:"password" validate:"required,min=8" gorm:"not null"`
	Activate bool    `json:"-" gorm:"default:false"`
}

type UserQueryable struct {
	Email    *string `json:"email" validate:"required,email" gorm:"not null"`
	Password *string `json:"password" validate:"required,min=8" gorm:"not null"`
}

type Users []User

func (User) TableName() string          { return "users" }
func (Users) TableName() string         { return User{}.TableName() }
func (UserCreatable) TableName() string { return User{}.TableName() }
func (UserUpdatable) TableName() string { return User{}.TableName() }
func (UserQueryable) TableName() string { return User{}.TableName() }
func (UserResponse) TableName() string  { return User{}.TableName() }

func (usr UserCreatable) Validate() error {
	validate := validator.New()
	if err := validate.Struct(&usr); err != nil {
		fmt.Println("Error while validate incoming request user: " + err.Error())
		return err
	}
	return nil
}

func (usr UserUpdatable) Validate() error {
	validate := validator.New()
	if err := validate.Struct(&usr); err != nil {
		fmt.Println("Error while validate incoming request user: " + err.Error())
		return err
	}
	return nil
}

func (usr UserQueryable) Validate() error {
	validate := validator.New()
	if err := validate.Struct(&usr); err != nil {
		fmt.Println("Error while validate incoming request user: " + err.Error())
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

func (usr *UserUpdatable) Mask() {
	hashPasswordBytes, _ := bcrypt.GenerateFromPassword([]byte(*usr.Password), 5)
	hashPassword := string(hashPasswordBytes)

	*usr.Password = hashPassword
}

func (usr UserCreatable) Convert2Response() UserResponse {
	return UserResponse{Model: usr.Model, UUID: usr.UUID, FirstName: *usr.FirstName, LastName: *usr.LastName, Gender: *usr.Gender, Email: *usr.Email, Password: *usr.Password, RoleId: *usr.RoleID, Activate: usr.Activate, ActivationCode: usr.ActivationCode, AvatarURL: usr.AvatarURL}
}

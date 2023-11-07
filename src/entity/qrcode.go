package entity

import (
	"fmt"
	"mime/multipart"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QRCode struct {
	gorm.Model `json:"-"`
	UserID     uint `json:"-"`

	UUID                  string `json:"uuid" gorm:"not null"`
	Content               string `json:"content" gorm:"not null"`
	Type                  string `json:"type" gorm:"default:text"`
	Background            string `json:"background" gorm:"default:#FFFFFF"`
	Foreground            string `json:"foreground" gorm:"default:#000000"`
	BorderWidth           int    `json:"borderWidth" gorm:"default:20"`
	CircleShape           bool   `json:"circleShape" gorm:"not null;default:false"`
	TransparentBackground bool   `json:"transparentBackground" gorm:"not null;default:false"`
	Version               int    `json:"version" gorm:"default:2"`
	ErrorLevel            int    `json:"errorLevel" gorm:"default:2"`
	PublicURL             string `json:"publicURL" gorm:"not null"`
	EncodeContent         string `json:"encode" gorm:"not null"`
	FilePath              string `json:"-" gorm:"not null"`
}

type QRCodeResponse struct {
	gorm.Model `json:"-"`

	UUID                  string `json:"uuid" gorm:"not null"`
	Content               string `json:"content" gorm:"not null"`
	Type                  string `json:"type" gorm:"default:text"`
	Background            string `json:"background" gorm:"default:#FFFFFF"`
	Foreground            string `json:"foreground" gorm:"default:#000000"`
	BorderWidth           int    `json:"borderWidth" gorm:"default:20"`
	CircleShape           bool   `json:"circleShape" gorm:"not null;default:false"`
	TransparentBackground bool   `json:"transparentBackground" gorm:"not null;default:false"`
	Version               int    `json:"version" gorm:"default:2"`
	ErrorLevel            int    `json:"errorLevel" gorm:"default:2"`
	PublicURL             string `json:"publicURL" gorm:"not null"`
	EncodeContent         string `json:"encodeContent" gorm:"not null"`
	FilePath              string `json:"-" gorm:"not null"`
}

type QRCodeCreatable struct {
	gorm.Model

	Content               *string `form:"content" json:"content" validate:"required" gorm:"not null"`
	Background            *string `form:"background" json:"background" validate:"hexcolor" gorm:"default:#FFFFFF"`
	Foreground            *string `form:"foreground" json:"foreground" validate:"hexcolor" gorm:"default:#000000"`
	BorderWidth           *int    `form:"borderWidth" json:"borderWidth" validate:"required" gorm:"default:20"`
	CircleShape           *bool   `form:"circleShape" json:"circleShape" gorm:"not null;default:false"`
	TransparentBackground *bool   `form:"transparentBackground" json:"transparentBackground" gorm:"not null;default:false"`
	ErrorLevel            *int    `form:"errorLevel" json:"errorLevel" validate:"gte=1,lte=4" gorm:"default:2"`
	Version               *int    `form:"-" json:"version" gorm:"default:2"`
	UserID                uint    `form:"-" json:"-" validate:"required" gorm:"not null"`
	UUID                  string  `form:"-" json:"-" gorm:"not null"`
	Type                  string  `form:"-" json:"-" gorm:"default:text"`
	PublicURL             string  `form:"-" json:"-" gorm:"not null"`
	EncodeContent         string  `form:"-" json:"-" gorm:"not null"`
	FilePath              string  `form:"-" json:"-" gorm:"not null"`

	Halftone *multipart.FileHeader `form:"halftone" sql:"-" gorm:"-"`
	Logo     *multipart.FileHeader `form:"logo" sql:"-" gorm:"-"`
}

type QRCodeUpdatable struct {
	gorm.Model
}

type QRCodes []QRCode

func (QRCode) TableName() string          { return "qrcodes" }
func (QRCodes) TableName() string         { return QRCode{}.TableName() }
func (QRCodeCreatable) TableName() string { return QRCode{}.TableName() }
func (QRCodeUpdatable) TableName() string { return QRCode{}.TableName() }

func (qrCode QRCodeCreatable) Validate() error {
	validate := validator.New()
	if err := validate.Struct(&qrCode); err != nil {
		fmt.Println("Error while validate incoming request qrcode: " + err.Error())
		return err
	}
	return nil
}

func (qrCode QRCodeUpdatable) Validate() error {
	validate := validator.New()
	if err := validate.Struct(&qrCode); err != nil {
		fmt.Println("Error while validate incoming request qrcode: " + err.Error())
		return err
	}
	return nil
}

func (qrCode *QRCodeCreatable) Mask() {
	qrCode.UUID = uuid.NewString()
}

package entity

import (
	"fmt"
	"go-qrcode-generator-cms-api/src/constants"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"gorm.io/gorm"
)

type QRCode struct {
	gorm.Model    `json:"-"`
	UserID        uint   `json:"-"`
	UUID          string `json:"-" gorm:"not null"`
	Content       string `json:"-" gorm:"not null"`
	Type          string `json:"-" gorm:"default:text"`
	Background    string `json:"-" gorm:"default:#FFFFFF"`
	Foreground    string `json:"-" gorm:"default:#000000"`
	BorderWidth   int    `json:"-" gorm:"default:20"`
	CircleShape   bool   `json:"-" gorm:"not null;default:false"`
	Version       int    `json:"-" gorm:"default:1"`
	ErrorLevel    int    `json:"-" gorm:"default:1"`
	PublicURL     string `json:"-" gorm:"column:public_url;not null"`
	EncodeContent string `json:"-" gorm:"column:encode_content;not null"`
	FilePath      string `json:"-" gorm:"column:file_path;not null"`
}

type QRCodeResponse struct {
	gorm.Model    `json:"-"`
	UserId        uint   `json:"-" gorm:"column:user_id"`
	User          Role   `json:"user"`
	UUID          string `json:"uuid" gorm:"not null"`
	Content       string `json:"content" gorm:"not null"`
	Type          string `json:"type" gorm:"default:text"`
	Background    string `json:"background" gorm:"default:#FFFFFF"`
	Foreground    string `json:"foreground" gorm:"default:#000000"`
	BorderWidth   int    `json:"borderWidth" gorm:"default:20"`
	CircleShape   bool   `json:"circleShape" gorm:"not null;default:false"`
	Version       int    `json:"version" gorm:"default:1"`
	ErrorLevel    int    `json:"errorLevel" gorm:"default:1"`
	PublicURL     string `json:"publicURL" gorm:"column:public_url;not null"`
	EncodeContent string `json:"encodeContent" gorm:"column:encode_content;not null"`
	FilePath      string `json:"-" gorm:"column:file_path;not null"`
}

// TODO: Custom error message
type QRCodeCreatable struct {
	gorm.Model    `json:"-"`
	UserID        *uint   `json:"-" validate:"required" gorm:"column:user_id;not null"`
	UUID          string  `json:"-" gorm:"not null"`
	Content       *string `json:"content" validate:"required" gorm:"not null"`
	Type          string  `json:"-" gorm:"default:text"`
	Background    *string `json:"background" validate:"hexcolor" gorm:"default:#FFFFFF"`
	Foreground    *string `json:"foreground" validate:"hexcolor" gorm:"default:#000000"`
	BorderWidth   *int    `json:"borderWidth" validate:"required" gorm:"default:20"`
	CircleShape   *bool   `json:"circleShape" gorm:"not null;default:false"`
	Version       *int    `json:"version" validate:"gte=1,lte=40" gorm:"default:1"`
	ErrorLevel    *int    `json:"errorLevel" validate:"gte=1,lte=4" gorm:"default:1"`
	PublicURL     string  `json:"-" gorm:"column:public_url;not null"`
	EncodeContent string  `json:"-" gorm:"column:encode_content;not null"`
	FilePath      string  `json:"-" gorm:"column:file_path;not null"`
}

type QRCodeUpdatable struct {
	gorm.Model `json:"-"`
}

type QRCodes []QRCode

func (QRCode) GetTableName() string          { return "qr_codes" }
func (QRCodes) GetTableName() string         { return Role{}.GetTableName() }
func (QRCodeCreatable) GetTableName() string { return QRCode{}.GetTableName() }
func (QRCodeUpdatable) GetTableName() string { return QRCode{}.GetTableName() }

func (qrCode QRCodeCreatable) Validate() error {
	validate := validator.New()
	if err := validate.Struct(&qrCode); err != nil {
		fmt.Println("Error while validate qrcode creatable: " + err.Error())
		return err
	}
	return nil
}

func (qrCode QRCodeUpdatable) Validate() error {
	validate := validator.New()
	if err := validate.Struct(&qrCode); err != nil {
		fmt.Println("Error while validate qrcode updatable: " + err.Error())
		return err
	}
	return nil
}

// TODO: Also check upper case use case
func DetectQRCodeType(content string) string {
	switch {
	case strings.Contains(content, "http") || strings.Contains(content, "https"):
		return constants.QRCodeURLType
	case strings.Contains(content, "mailto"):
		return constants.QRCodeMailType
	case strings.Contains(content, "smsto"):
		return constants.QRCodeSMSType
	case strings.Contains(content, "tel"):
		return constants.QRCodeTelType
	case strings.Contains(content, "wifi"):
		return constants.QRCodeWifiType
	default:
		return constants.QRCodeTextType
	}
}

// TODO: Optimize this function later
func (qrCode *QRCodeCreatable) Mask() {
	currentDir, _ := os.Getwd()
	basePath := os.Getenv("LOCAL_QRCODE_STORAGE_DIR")
	defaultFileType := os.Getenv("LOCAL_QRCODE_FILE_TYPE")
	fileName := uuid.NewString() + "." + defaultFileType
	filePath := filepath.Join(currentDir, basePath, fileName)

	qrCode.UUID = uuid.NewString()
	qrCode.Type = DetectQRCodeType(*qrCode.Content)
	qrCode.FilePath = filePath
}

func (qrCode QRCodeCreatable) Standardized() (*[]qrcode.EncodeOption, *[]standard.ImageOption) {
	qrCodeConfigs := []qrcode.EncodeOption{
		qrcode.WithEncodingMode(qrcode.EncModeByte),
		qrcode.WithVersion(*qrCode.Version),
	}
	writerConfigs := []standard.ImageOption{
		standard.WithBorderWidth(*qrCode.BorderWidth),
	}

	if *qrCode.ErrorLevel == 1 {
		qrCodeConfigs = append(qrCodeConfigs, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionLow))
	} else if *qrCode.ErrorLevel == 2 {
		qrCodeConfigs = append(qrCodeConfigs, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionMedium))
	} else if *qrCode.ErrorLevel == 3 {
		qrCodeConfigs = append(qrCodeConfigs, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionQuart))
	} else {
		qrCodeConfigs = append(qrCodeConfigs, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest))
	}
	if qrCode.CircleShape != nil && *qrCode.CircleShape {
		writerConfigs = append(writerConfigs, standard.WithCircleShape())
	}
	return &qrCodeConfigs, &writerConfigs
}

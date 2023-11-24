package utils

import (
	"errors"
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"
	"mime/multipart"
	"os"
	"strconv"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

var (
	ErrCreateQREncodeOption    = errors.New("create QRCode encode option failure")
	ErrCreateQRImageOption     = errors.New("create QRCode image option failure")
	ErrSaveQRCode2LocalStorage = errors.New("save QrCode to local storage failure")
	ErrQRCodeLogoSizeTooLarge  = errors.New("verify QrCode logo size failure: file size too large")
)

type qrCodeUtil struct {
}

func NewQrCodeUtil() *qrCodeUtil {
	return &qrCodeUtil{}
}

func (qrCodeUtil) SaveQRCode2LocalStorage(qrCode *entity.QRCodeCreatable, qrCodeConfigs []qrcode.EncodeOption, writerConfigs []standard.ImageOption) (*string, error) {
	if qrc, err := qrcode.NewWith(*qrCode.Content, qrCodeConfigs...); err != nil {
		fmt.Println("Error while create QRCode encode option: " + err.Error())
		return nil, ErrCreateQREncodeOption
	} else if w, err := standard.New(qrCode.FilePath, writerConfigs...); err != nil {
		fmt.Println("Error while create QRCode image option: " + err.Error())
		return nil, ErrCreateQRImageOption
	} else if err := qrc.Save(w); err != nil {
		fmt.Println("Error while save QrCode to local storage: " + err.Error())
		return nil, ErrSaveQRCode2LocalStorage
	}
	return &qrCode.FilePath, nil
}

func (qrCodeUtil) VerifyQrCodeLogoSize(logo *multipart.FileHeader) error {
	logoSizeInByte := logo.Size
	maxLogoSizeInByte, _ := strconv.Atoi(os.Getenv("MAX_QRLOGO_SIZE_IN_BYTE"))
	fmt.Println("QRCode Logo file size in bytes: ", logoSizeInByte)
	if logoSizeInByte > int64(maxLogoSizeInByte) {
		return ErrQRCodeLogoSizeTooLarge
	}
	return nil
}

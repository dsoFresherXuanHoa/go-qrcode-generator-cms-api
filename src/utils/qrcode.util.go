package utils

import (
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

type qrCodeUtil struct {
}

func NewQrCodeUtil() *qrCodeUtil {
	return &qrCodeUtil{}
}

func (qrCodeUtil) SaveQRCode2LocalStorage(qrCode *entity.QRCodeCreatable) (*string, error) {
	qrCodeConfigs, writerConfigs := qrCode.Standardized()
	if qrc, err := qrcode.NewWith(*qrCode.Content, *qrCodeConfigs...); err != nil {
		fmt.Println("Error while create qrcode.QrCode struct: " + err.Error())
		return nil, err
	} else if w, err := standard.New(qrCode.FilePath, *writerConfigs...); err != nil {
		fmt.Println("Error while create qrcode writer struct: " + err.Error())
		return nil, err
	} else if err := qrc.Save(w); err != nil {
		fmt.Println("Error while save QrCode to local server: " + err.Error())
		return nil, err
	}
	return &qrCode.FilePath, nil
}

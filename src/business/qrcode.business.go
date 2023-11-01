package business

import (
	"context"
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"
	"go-qrcode-generator-cms-api/src/storage"
	"go-qrcode-generator-cms-api/src/utils"

	"github.com/cloudinary/cloudinary-go/v2"
)

type QRCodeStorage interface {
	CreateQRCode(ctx context.Context, qrCode *entity.QRCodeCreatable) (*string, error)
}

type qrCodeBusiness struct {
	qrCodeStorage QRCodeStorage
}

func NewQRCodeBusiness(qrCodeStorage QRCodeStorage) *qrCodeBusiness {
	return &qrCodeBusiness{qrCodeStorage: qrCodeStorage}
}

func (business *qrCodeBusiness) CreateQRCode(ctx context.Context, cld *cloudinary.Cloudinary, qrCode *entity.QRCodeCreatable) (*string, *string, error) {
	if err := qrCode.Validate(); err != nil {
		fmt.Println("Error while validate qrcode request in qrcode business: " + err.Error())
		return nil, nil, err
	} else {
		qrCode.Mask()
		if filePath, err := utils.NewQrCodeUtil().SaveQRCode2LocalStorage(qrCode); err != nil {
			fmt.Println("Error while save QrCode to local server: " + err.Error())
			return nil, nil, err
		} else if encode, err := utils.NewImageUtil().Image2Base64(*filePath); err != nil {
			fmt.Println("Error while encode QrCode image content to base64 format: " + err.Error())
			return filePath, nil, err
		} else {
			qrCode.EncodeContent = *encode
			cloudinaryStorage := storage.NewCloudinaryStore(cld)
			if uploadResult, err := cloudinaryStorage.UploadSingleImage(ctx, qrCode.EncodeContent); err != nil {
				fmt.Println("Error while upload single image to Cloudinary API: " + err.Error())
				return filePath, nil, err
			} else {
				qrCode.PublicURL = uploadResult.URL
				if _, err := business.qrCodeStorage.CreateQRCode(ctx, qrCode); err != nil {
					fmt.Println("Error while save qrCode information to database in qrcode business: " + err.Error())
					return filePath, nil, err
				}
				return encode, &uploadResult.URL, nil
			}
		}
	}
}

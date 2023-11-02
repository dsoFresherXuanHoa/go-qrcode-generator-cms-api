package storage

import (
	"context"
	"errors"
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"
)

var (
	ErrSaveQRCode = errors.New("save qrcode into database failure")
)

type qrCodeStorage struct {
	sql *sqlStorage
}

func NewQrCodeStore(sql *sqlStorage) *qrCodeStorage {
	return &qrCodeStorage{sql: sql}
}

func (s *qrCodeStorage) CreateQRCode(ctx context.Context, qrCode *entity.QRCodeCreatable) (*string, error) {
	if err := s.sql.db.Table(entity.QRCodeCreatable{}.TableName()).Create(&qrCode).Error; err != nil {
		fmt.Println("Error while save qrcode into database: " + err.Error())
		return nil, ErrSaveQRCode
	}
	return &qrCode.UUID, nil
}

package storage

import (
	"context"
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"
)

type qrCodeStorage struct {
	sql *sqlStorage
}

func NewQrCodeStore(sql *sqlStorage) *qrCodeStorage {
	return &qrCodeStorage{sql: sql}
}

func (s *qrCodeStorage) CreateQRCode(ctx context.Context, qrCode *entity.QRCodeCreatable) (*string, error) {
	if err := s.sql.db.Table(entity.QRCodeCreatable{}.GetTableName()).Create(&qrCode).Error; err != nil {
		fmt.Println("Error while save qrCode information to database in qrCode storage: " + err.Error())
		return nil, err
	}
	return &qrCode.UUID, nil
}

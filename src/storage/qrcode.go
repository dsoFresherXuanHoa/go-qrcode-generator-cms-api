package storage

import (
	"context"
	"errors"
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"

	"github.com/redis/go-redis/v9"
)

var (
	ErrSaveQRCode       = errors.New("save qrcode into database failure")
	ErrFindQRCodeByUUID = errors.New("find qrcode by uuid failure")
)

type qrCodeStorage struct {
	sql   *sqlStorage
	redis *redisStorage
}

func NewQrCodeStore(sql *sqlStorage, redis *redisStorage) *qrCodeStorage {
	return &qrCodeStorage{sql: sql, redis: redis}
}

func (s *qrCodeStorage) CreateQRCode(ctx context.Context, client *redis.Client, qrCode *entity.QRCodeCreatable) (*string, error) {
	if err := s.sql.db.Table(entity.QRCodeCreatable{}.TableName()).Create(&qrCode).Error; err != nil {
		fmt.Println("Error while save qrcode into database: " + err.Error())
		return nil, ErrSaveQRCode
	} else if _, err := s.redis.SaveQRCode(client, qrCode); err != nil {
		return nil, err
	}
	return &qrCode.UUID, nil
}

func (s *qrCodeStorage) FindQRCodeByUUID(ctx context.Context, uuid string) (*entity.QRCodeResponse, error) {
	var qrCode entity.QRCode
	if err := s.sql.db.Where("uuid = ?", uuid).First(&qrCode).Error; err != nil {
		fmt.Println("Error while find qrcode by uuid: " + err.Error())
		return nil, ErrFindQRCodeByUUID
	}
	res := qrCode.Convert2Response()
	return &res, nil
}

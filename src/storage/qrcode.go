package storage

import (
	"context"
	"errors"
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"

	"github.com/redis/go-redis/v9"
)

var (
	ErrSaveQRCode = errors.New("save qrcode into database failure")
)

type qrCodeStorage struct {
	sql   *sqlStorage
	redis *redisStorage
}

func NewQrCodeStore(sql *sqlStorage, redis *redisStorage) *qrCodeStorage {
	return &qrCodeStorage{sql: sql, redis: redis}
}

func (s *qrCodeStorage) CreateQRCode(ctx context.Context, client *redis.Client, qrCode *entity.QRCodeCreatable) (*string, error) {
	if _, err := s.redis.SaveQRCode(client, qrCode); err != nil {
		return nil, err
	} else if err := s.sql.db.Table(entity.QRCodeCreatable{}.TableName()).Create(&qrCode).Error; err != nil {
		fmt.Println("Error while save qrcode into database: " + err.Error())
		return nil, ErrSaveQRCode
	}
	return &qrCode.UUID, nil
}

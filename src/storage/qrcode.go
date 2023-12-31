package storage

import (
	"context"
	"errors"
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"
	"strconv"
	"time"
)

var (
	ErrSaveQRCode            = errors.New("save qrcode into database failure")
	ErrFindQRCodeByUUID      = errors.New("find qrcode by uuid failure")
	ErrFindQRCodeByCondition = errors.New("find qrcode by condition failure")
	ErrFindAllQRCode         = errors.New("find all qrcode by condition failure")
)

type qrCodeStorage struct {
	sql   *sqlStorage
	redis *redisStorage
}

func NewQrCodeStore(sql *sqlStorage, redis *redisStorage) *qrCodeStorage {
	return &qrCodeStorage{sql: sql, redis: redis}
}

func (s *qrCodeStorage) CreateQRCode(ctx context.Context, qrCode *entity.QRCodeCreatable) (*string, error) {
	if err := s.sql.db.Table(entity.QRCodeCreatable{}.TableName()).Create(&qrCode).Error; err != nil {
		fmt.Println("Error while save qrcode into database: " + err.Error())
		return nil, ErrSaveQRCode
	} else if _, err := s.redis.SaveQRCode(qrCode); err != nil {
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

func (s *qrCodeStorage) FindQRCodeByCondition(ctx context.Context, cond map[string]interface{}, timeStat map[string]string, paging *entity.Paginate) ([]entity.QRCodeResponse, error) {
	var qrCodes entity.QRCodes
	offset := (paging.Page - 1) * paging.Size
	limit := paging.Size
	startTimeUnix, _ := strconv.ParseInt(timeStat["start_time"], 10, 64)
	endTimeUnix, _ := strconv.ParseInt(timeStat["end_time"], 10, 64)
	startTime := time.Unix(startTimeUnix, 0)
	endTime := time.Unix(endTimeUnix, 0)

	var total int64
	if endTime.After(startTime) {
		if err := s.sql.db.Where("created_at > ? AND created_at < ?", startTime, endTime).Where(cond).Offset(offset).Limit(limit).Find(&qrCodes).Error; err != nil {
			s.sql.db.Table(entity.QRCode{}.TableName()).Count(&total)
			fmt.Println("Error while find qrcode by condition with time filter: " + err.Error())
			return nil, ErrFindQRCodeByCondition
		}
	} else if err := s.sql.db.Where(cond).Offset(offset).Limit(limit).Find(&qrCodes).Error; err != nil {
		s.sql.db.Table(entity.QRCode{}.TableName()).Count(&total)
		fmt.Println("Error while find qrcode by condition: " + err.Error())
		return nil, ErrFindQRCodeByCondition
	}

	paging.Total = int(total)
	var res = make([]entity.QRCodeResponse, len(qrCodes))
	for i, qrCode := range qrCodes {
		res[i] = qrCode.Convert2Response()
	}
	return res, nil
}

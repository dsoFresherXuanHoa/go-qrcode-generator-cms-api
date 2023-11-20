package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrGetQRCodeFromRedis         = errors.New("get QR Code content from Redis failure")
	ErrQRCodeRedisKeyNotFound     = errors.New("get QR Code by key failure: key not found")
	ErrSaveQRCode2Redis           = errors.New("save QR Code to Redis failure")
	ErrSaveAccessToken2Redis      = errors.New("save accessToken to Redis failure")
	ErrDeleteAccessTokenFromRedis = errors.New("delete accessToken from Redis failure")
	ErrGetAccessTokenFromRedis    = errors.New("get accessToken from Redis failure")
)

type redisStorage struct {
	redisClient *redis.Client
}

func NewRedisStore(redisClient *redis.Client) *redisStorage {
	return &redisStorage{redisClient: redisClient}
}

func (s *redisStorage) GetQRCodeEncodeFromRedis(key string) ([]string, error) {
	if qrCodeBase64 := s.redisClient.Get(context.Background(), key); qrCodeBase64.Err() != nil {
		fmt.Println("Error while get QR Code content from Redis: " + qrCodeBase64.Err().Error())
		return nil, ErrGetQRCodeFromRedis
	} else if result, err := qrCodeBase64.Result(); err == redis.Nil {
		fmt.Println("Error key not found in Redis: " + err.Error())
		return nil, ErrQRCodeRedisKeyNotFound
	} else {
		res := []string{}
		json.Unmarshal([]byte(result), &res)
		return res, nil
	}
}

func (s *redisStorage) GetRedisKey(qrCode *entity.QRCodeCreatable) string {
	var result string
	if qrCode.Content != nil {
		result += *qrCode.Content
	}
	if qrCode.Background != nil {
		result += *qrCode.Background
	} else {
		result += "#FFFFFF"
	}
	if qrCode.Foreground != nil {
		result += *qrCode.Foreground
	} else {
		result += "#000000"
	}
	if qrCode.BorderWidth != nil {
		result += strconv.Itoa(*qrCode.BorderWidth)
	} else {
		result += "20"
	}
	if qrCode.CircleShape != nil && *qrCode.CircleShape {
		result += "true"
	} else {
		result += "false"
	}
	if qrCode.ErrorLevel != nil {
		result += strconv.Itoa(*qrCode.ErrorLevel)
	} else {
		result += "2"
	}
	if qrCode.Logo != nil {
		logoSize := qrCode.Logo.Size
		result += qrCode.Logo.Filename
		result += strconv.Itoa(int(logoSize))
	}
	return result
}

func (s *redisStorage) SaveQRCode(qrCode *entity.QRCodeCreatable) (*string, error) {
	key := s.GetRedisKey(qrCode)
	encode := &qrCode.EncodeContent
	publicUrl := qrCode.PublicURL
	value, _ := json.Marshal([]string{*encode, publicUrl})
	expireDuration, _ := strconv.Atoi(os.Getenv("REDIS_EXPIRE_TIME_IN_MINUTE"))
	if _, err := s.redisClient.Set(context.Background(), key, value, time.Minute*time.Duration(expireDuration)).Result(); err != nil {
		fmt.Println("Error while save QRCode to Redis Server: " + err.Error())
		return nil, ErrSaveQRCode2Redis
	}
	res := string(value)
	fmt.Println("Save QRCode to Redis Server Success!")
	return &res, nil
}

func (s *redisStorage) SaveAccessToken(key string, accessToken string) error {
	expireDuration, _ := strconv.Atoi(os.Getenv("REDIS_EXPIRE_TIME_IN_MINUTE"))
	if _, err := s.redisClient.Set(context.Background(), key, accessToken, time.Minute*time.Duration(expireDuration)).Result(); err != nil {
		fmt.Println("Error while save accessToken to Redis Server: " + err.Error())
		return ErrSaveAccessToken2Redis
	}
	return nil
}

func (s *redisStorage) DeleteAccessToken(key string) error {
	if _, err := s.redisClient.Del(context.Background(), key).Result(); err != nil {
		fmt.Println("Error while delete accessToken from Redis Server: " + err.Error())
		return ErrDeleteAccessTokenFromRedis
	}
	return nil
}

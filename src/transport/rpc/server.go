package rpc

import (
	"context"
	"errors"
	"fmt"
	"go-qrcode-generator-cms-api/proto/gRPC/qrcodes"
	"go-qrcode-generator-cms-api/src/business"
	"go-qrcode-generator-cms-api/src/storage"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ErrCreateQRCodeViaGRPCService = errors.New("create qrCode via GRPC Service failure")
)

type server struct {
	db          *gorm.DB
	redisClient *redis.Client
	qrcodes.UnimplementedQRCodeServiceServer
}

func NewGrpcServer(db *gorm.DB, redisClient *redis.Client) *server {
	return &server{db: db, redisClient: redisClient}
}

func (s server) GrpcCreateQRCode(ctx context.Context, req *qrcodes.CreateQRCodeRequest) (*qrcodes.CreateQRCodeResponse, error) {
	sqlStorage := storage.NewSQLStore(s.db)
	redisStorage := storage.NewRedisStore(s.redisClient)
	qrCodeStorage := storage.NewQrCodeStore(sqlStorage, redisStorage)
	qrCodeBusiness := business.NewQRCodeBusiness(qrCodeStorage, redisStorage)
	if res, err := qrCodeBusiness.GrpcCreateQRCode(ctx, req); err != nil {
		fmt.Println("Error while create QRCode via GRPC Service: " + err.Error())
		return nil, ErrCreateQRCodeViaGRPCService
	} else {
		return res, nil
	}
}

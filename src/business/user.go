package business

import (
	"context"
	"go-qrcode-generator-cms-api/src/entity"
)

var ()

type UserStorage interface {
	FindQRCodeByUserId(ctx context.Context, userUUID string, cond map[string]interface{}, timeStat map[string]string, paging entity.Paginate) ([]entity.QRCodeResponse, error)
}

type userBusiness struct {
	userStorage UserStorage
}

func NewUserBusiness(userStorage UserStorage) *userBusiness {
	return &userBusiness{userStorage: userStorage}
}

func (business *userBusiness) FindQRCodeByUserId(ctx context.Context, userUUID string, cond map[string]interface{}, timeStat map[string]string, paging *entity.Paginate) ([]entity.QRCodeResponse, error) {
	paging.Standardized()
	if qrCodes, err := business.userStorage.FindQRCodeByUserId(ctx, userUUID, cond, timeStat, *paging); err != nil {
		return nil, err
	} else {
		return qrCodes, nil
	}
}

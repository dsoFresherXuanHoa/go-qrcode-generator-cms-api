package storage

import (
	"context"
	"errors"
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrSaveUser                                 = errors.New("save user failure")
	ErrUpdateUserActivateStatusByActivationCode = errors.New("update user activate status by activation code failure")
	ErrFindUserByEmailAndActivationCode         = errors.New("find user by email and activate status failure")
	ErrCompareUserPassword                      = errors.New("compare user password and hash password failure")
	ErrFindUserById                             = errors.New("find detail user by id failure")
	ErrFindUserByActivationCode                 = errors.New("find user by activation code failure")
	ErrUpdateUserPasswordByActivationCode       = errors.New("update user password by activation code failure")
	ErrFindUserByEmail                          = errors.New("find user by email failure")
	ErrFindQRCodeByUserId                       = errors.New("find qr code by userId failure")
)

type userStorage struct {
	sql *sqlStorage
}

func NewUserStore(sql *sqlStorage) *userStorage {
	return &userStorage{sql: sql}
}

func (s *userStorage) CreateUser(ctx context.Context, user *entity.UserCreatable) (*string, error) {
	if err := s.sql.db.Table(entity.UserCreatable{}.TableName()).Create(&user).Error; err != nil {
		fmt.Println("Error while save user into database: " + err.Error())
		return nil, ErrSaveRole2DB
	}
	return &user.UUID, nil
}

func (s *userStorage) UpdateUserActivateStatusByActivationCode(ctx context.Context, activationCode string) error {
	if err := s.sql.db.Model(&entity.User{}).Where("activation_code = ?", activationCode).Update("activate", true).Error; err != nil {
		fmt.Println("Error while update user activate status by activation code into database: " + err.Error())
		return ErrUpdateUserActivateStatusByActivationCode
	}
	return nil
}

func (s *userStorage) FindUserByEmailAndPassword(ctx context.Context, email string, password string) (*entity.UserResponse, error) {
	var usr entity.UserResponse
	if err := s.sql.db.Table(usr.TableName()).Where("email = ?", email).Where("activate = ?", true).First(&usr).Error; err != nil {
		fmt.Println("Error while get user by email and activation status from database: " + err.Error())
		return nil, ErrFindUserByEmailAndActivationCode
	} else if err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password)); err != nil {
		fmt.Println("Error while compare user password and hash password from database: " + err.Error())
		return nil, ErrCompareUserPassword
	}
	return &usr, nil
}

func (s *userStorage) FindDetailUserById(ctx context.Context, id uint) (*entity.UserResponse, error) {
	var usr entity.UserResponse
	if err := s.sql.db.Table(usr.TableName()).Where("id = ?", id).Preload("Role").First(&usr).Error; err != nil {
		fmt.Println("Error while find detail user by id from database: " + err.Error())
		return nil, ErrFindUserById
	}
	return &usr, nil
}

func (s *userStorage) FindUserByActivationCode(ctx context.Context, activationCode string) (*entity.User, error) {
	var usr entity.User
	if err := s.sql.db.Table(usr.TableName()).Where("activation_code = ?", activationCode).First(&usr).Error; err != nil {
		fmt.Println("Error while find user by activationCode from database: " + err.Error())
		return nil, ErrFindUserByActivationCode
	}
	return &usr, nil
}

func (s *userStorage) UpdateUserPasswordByActivationCode(ctx context.Context, activationCode string, user *entity.UserUpdatable) (*string, error) {
	if usr, err := s.FindUserByActivationCode(ctx, activationCode); err != nil {
		return nil, err
	} else if err := s.sql.db.Model(&user).Where("activation_code = ?", activationCode).Update("password", *user.Password).Error; err != nil {
		fmt.Println("Error while update user password by activation code into database: " + err.Error())
		return nil, ErrUpdateUserPasswordByActivationCode
	} else {
		userId := fmt.Sprint(usr.ID)
		return &userId, nil
	}
}

func (s *userStorage) FindUserByEmail(ctx context.Context, email string) (*entity.UserCreatable, error) {
	var usr entity.UserCreatable
	if err := s.sql.db.Table(usr.TableName()).Where("email = ?", email).Where("activate = ?", true).First(&usr).Error; err != nil {
		fmt.Println("Error while find user by email from database: " + err.Error())
		return nil, ErrFindUserByEmail
	}
	return &usr, nil
}

func (s *userStorage) FindQRCodeByUserId(ctx context.Context, userUUID string, cond map[string]interface{}, timeStat map[string]string, paging entity.Paginate) ([]entity.QRCodeResponse, error) {
	var qrCodes entity.QRCodes
	offset := (paging.Page - 1) * paging.Size
	limit := paging.Size
	startTimeUnix, _ := strconv.ParseInt(timeStat["start_time"], 10, 64)
	endTimeUnix, _ := strconv.ParseInt(timeStat["end_time"], 10, 64)
	startTime := time.Unix(startTimeUnix, 0)
	endTime := time.Unix(endTimeUnix, 0)
	fmt.Println(startTime, endTime)
	if endTime.After(startTime) {
		if err := s.sql.db.Where("created_at > ? AND created_at < ?", startTime, endTime).Where(cond).Offset(offset).Limit(limit).Find(&qrCodes).Error; err != nil {
			fmt.Println("Error while find qrcode by condition with time filter: " + err.Error())
			return nil, ErrFindQRCodeByCondition
		}
	} else if err := s.sql.db.Where(cond).Offset(offset).Limit(limit).Find(&qrCodes).Error; err != nil {
		fmt.Println("Error while find qrcode by condition: " + err.Error())
		return nil, ErrFindQRCodeByCondition
	}
	var res = make([]entity.QRCodeResponse, len(qrCodes))
	for i, qrCode := range qrCodes {
		res[i] = qrCode.Convert2Response()
	}
	return res, nil
}

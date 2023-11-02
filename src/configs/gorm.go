package configs

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrLoadGormEnvFile = errors.New("load .env file failure")
	ErrConnect2DB      = errors.New("connect to database via GORM failure")
)

type gormClient struct {
	instance *gorm.DB
}

func NewGormClient() *gormClient {
	return &gormClient{instance: nil}
}

func (instance *gormClient) Instance() (*gorm.DB, error) {
	if instance.instance == nil {
		if err := godotenv.Load(); err != nil {
			fmt.Println("Error while load .env file: " + err.Error())
			return nil, ErrLoadGormEnvFile
		} else {
			var dns = os.Getenv("GORM_URL")
			if database, err := gorm.Open(mysql.Open(dns), &gorm.Config{}); err != nil {
				fmt.Println("Error while connect to database via GORM: " + err.Error())
				return nil, ErrConnect2DB
			} else {
				instance.instance = database
			}
		}
	}
	return instance.instance, nil
}

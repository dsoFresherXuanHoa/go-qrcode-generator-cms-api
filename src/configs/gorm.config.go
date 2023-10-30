package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type gormInstance struct {
	instance *gorm.DB
}

func NewGormInstance() *gormInstance {
	return &gormInstance{instance: nil}
}

func (instance *gormInstance) GetGormInstance() (*gorm.DB, error) {
	if instance.instance == nil {
		if err := godotenv.Load(); err != nil {
			fmt.Println("Can't load .env files! Check your .env file and try again later: " + err.Error())
			return nil, err
		} else {
			var dns = os.Getenv("GORM_URL")
			if database, err := gorm.Open(mysql.Open(dns), &gorm.Config{}); err != nil {
				fmt.Println("Can't connect to database using GORM: " + err.Error())
				return nil, err
			} else {
				instance.instance = database
			}
		}
	}
	fmt.Println("Connection to database has been created!!!")
	return instance.instance, nil
}

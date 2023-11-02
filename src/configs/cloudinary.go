package configs

import (
	"errors"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/joho/godotenv"
)

var (
	ErrLoadCloudinaryEnvFile = errors.New("load .env file failure")
	ErrConnect2Cloudinary    = errors.New("connect to cloudinary failure")
)

type cloudinaryClient struct {
	instance *cloudinary.Cloudinary
}

func NewCloudinaryClient() *cloudinaryClient {
	return &cloudinaryClient{instance: nil}
}

func (instance *cloudinaryClient) Instance() (*cloudinary.Cloudinary, error) {
	if instance.instance == nil {
		if err := godotenv.Load(); err != nil {
			fmt.Println("Error while load .env file: " + err.Error())
			return nil, ErrLoadCloudinaryEnvFile
		} else {
			var url = os.Getenv("CLOUDINARY_URL")
			if cld, err := cloudinary.NewFromURL(url); err != nil {
				fmt.Println("Error while connect to Cloudinary: " + err.Error())
				return nil, ErrConnect2Cloudinary
			} else {
				instance.instance = cld
			}
		}
	}
	return instance.instance, nil
}

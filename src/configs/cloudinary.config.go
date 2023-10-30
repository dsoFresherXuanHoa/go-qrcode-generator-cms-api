package configs

import (
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/joho/godotenv"
)

type cloudinaryInstance struct {
	instance *cloudinary.Cloudinary
}

func NewCloudinaryInstance() *cloudinaryInstance {
	return &cloudinaryInstance{instance: nil}
}

func (instance *cloudinaryInstance) GetCloudinaryInstance() (*cloudinary.Cloudinary, error) {
	if instance.instance == nil {
		if err := godotenv.Load(); err != nil {
			fmt.Println("Can't load .env files! Check your .env file and try again later: " + err.Error())
			return nil, err
		} else {
			var url = os.Getenv("CLOUDINARY_URL")
			if cld, err := cloudinary.NewFromURL(url); err != nil {
				fmt.Println("Can't connect to cloudinary using cloudinary API: " + err.Error())
				return nil, err
			} else {
				fmt.Println("Connection to cloudinary has been created!!!")
				instance.instance = cld
			}
		}
	}
	return instance.instance, nil
}

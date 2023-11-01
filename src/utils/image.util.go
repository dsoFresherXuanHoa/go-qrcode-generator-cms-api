package utils

import (
	"encoding/base64"
	"fmt"
	"go-qrcode-generator-cms-api/src/errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"golang.org/x/exp/slices"
)

type imageUtil struct{}

func NewImageUtil() *imageUtil {
	return &imageUtil{}
}

func (imageUtil) ImageFileHeader2Base64(file *multipart.FileHeader) (*string, error) {
	fileContent, _ := file.Open()
	var validImageType = []string{"image/jpeg", "image/png", "image/jpg", "image/webp"}
	if bytes, err := io.ReadAll(fileContent); err != nil {
		fmt.Println("Error while encoding image from request header to base64 format: " + err.Error())
		return nil, err
	} else if mimeType := http.DetectContentType(bytes); !slices.Contains(validImageType, mimeType) {
		fmt.Println("Error while upload image to Cloudinary: invalid image type")
		return nil, errors.ErrInvalidImageType
	} else {
		result := "data:" + mimeType + ";base64,"
		result += base64.StdEncoding.EncodeToString(bytes)
		return &result, nil
	}
}

func (imageUtil) Image2Base64(filePath string) (*string, error) {
	var validImageType = []string{"image/jpeg", "image/png", "image/jpg", "image/webp"}
	if bytes, err := os.ReadFile(filePath); err != nil {
		fmt.Println("Error while encoding image from local file path to base64 format: " + err.Error())
		return nil, err
	} else if mimeType := http.DetectContentType(bytes); !slices.Contains(validImageType, mimeType) {
		fmt.Println("Error while validate image format: invalid image type")
		return nil, errors.ErrInvalidImageType
	} else {
		result := "data:" + mimeType + ";base64,"
		result += base64.StdEncoding.EncodeToString(bytes)
		return &result, nil
	}
}

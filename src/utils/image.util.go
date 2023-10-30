package utils

import (
	"encoding/base64"
	"fmt"
	"go-qrcode-generator-cms-api/src/exception"
	"io"
	"mime/multipart"
	"net/http"

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
		return nil, exception.ErrInvalidImageType
	} else {
		result := "data:" + mimeType + ";base64,"
		result += base64.StdEncoding.EncodeToString(bytes)
		return &result, nil
	}
}

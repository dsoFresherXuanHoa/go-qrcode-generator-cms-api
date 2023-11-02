package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"golang.org/x/exp/slices"
)

var (
	ErrReadByteFromImagePath           = errors.New("read all byte from image path failure")
	ErrUnsupportedImageType            = errors.New("validate image type failure: unsupported image type")
	ErrReadByteFromMultipartFileHeader = errors.New("read all byte from image path failure")
)

type imageUtil struct{}

func NewImageUtil() *imageUtil {
	return &imageUtil{}
}

func (imageUtil) ImageMultipartFile2Base64(file *multipart.FileHeader) (*string, error) {
	fileContent, _ := file.Open()
	var supportedMineType = []string{"image/jpeg", "image/png", "image/jpg"}
	if imageBytes, err := io.ReadAll(fileContent); err != nil {
		fmt.Println("Error while read all bytes from multipart file header: " + err.Error())
		return nil, ErrReadByteFromMultipartFileHeader
	} else if mimeType := http.DetectContentType(imageBytes); !slices.Contains(supportedMineType, mimeType) {
		fmt.Println("Error while validate image format: unsupported image type")
		return nil, ErrUnsupportedImageType
	} else {
		result := "data:" + mimeType + ";base64,"
		result += base64.StdEncoding.EncodeToString(imageBytes)
		return &result, nil
	}
}

func (imageUtil) Image2Base64(filePath string) (*string, error) {
	var supportedMineType = []string{"image/jpeg", "image/png", "image/jpg"}
	if imageBytes, err := os.ReadFile(filePath); err != nil {
		fmt.Println("Error while read all byte from image: " + err.Error())
		return nil, ErrReadByteFromImagePath
	} else if mimeType := http.DetectContentType(imageBytes); !slices.Contains(supportedMineType, mimeType) {
		fmt.Println("Error while validate image format: unsupported image type")
		return nil, ErrUnsupportedImageType
	} else {
		result := "data:" + mimeType + ";base64,"
		result += base64.StdEncoding.EncodeToString(imageBytes)
		return &result, nil
	}
}

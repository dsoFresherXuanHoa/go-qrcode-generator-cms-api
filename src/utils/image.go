package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

var (
	ErrReadByteFromImagePath           = errors.New("read all byte from image path failure")
	ErrUnsupportedImageType            = errors.New("validate image type failure: unsupported image type")
	ErrReadByteFromMultipartFileHeader = errors.New("read all byte from image file header failure")
	ErrOpenMultipartFileHeader         = errors.New("open image from image file header failure")
	ErrCreateEmptyLocalHalftone        = errors.New("create empty local halftone failure")
	ErrCopyFileHeader2LocalHalftone    = errors.New("copy image file header to local halftone failure")
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

func (imageUtil) ImageMultipartFile2LocalStorage(file *multipart.FileHeader) (*string, error) {
	currentDir, _ := os.Getwd()
	defaultLocalFileType := os.Getenv("LOCAL_QRCODE_HALFTONE_FILE_TYPE")
	localStorageBasePath := os.Getenv("LOCAL_QRCODE_HALFTONE_STORAGE_DIR")
	localFileName := uuid.NewString() + "." + defaultLocalFileType
	localFilePath := filepath.Join(currentDir, localStorageBasePath, localFileName)
	var supportedMineType = []string{"image/png"}
	if halftone, err := file.Open(); err != nil {
		fmt.Println("Error while open multipart file: " + err.Error())
		return nil, ErrOpenMultipartFileHeader
	} else if halftoneBytes, err := io.ReadAll(halftone); err != nil {
		fmt.Println("Error while read all bytes from multipart file header: " + err.Error())
		return nil, ErrReadByteFromMultipartFileHeader
	} else if mimeType := http.DetectContentType(halftoneBytes); !slices.Contains(supportedMineType, mimeType) {
		fmt.Println("Error while validate image format: unsupported image type")
		return nil, ErrUnsupportedImageType
	} else if halftoneFile, err := os.Create(localFilePath); err != nil {
		fmt.Println("Error while create empty local halftone file: " + err.Error())
		return nil, ErrCreateEmptyLocalHalftone
	} else if _, err := io.Copy(halftoneFile, bytes.NewReader(halftoneBytes)); err != nil {
		fmt.Println("Error while copy multipart file to local halftone file: " + err.Error())
		return nil, ErrCopyFileHeader2LocalHalftone
	} else {
		return &localFilePath, nil
	}
}

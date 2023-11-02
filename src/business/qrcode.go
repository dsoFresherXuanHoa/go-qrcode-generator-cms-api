package business

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"go-qrcode-generator-cms-api/src/constants"
	"go-qrcode-generator-cms-api/src/entity"
	"go-qrcode-generator-cms-api/src/storage"
	"go-qrcode-generator-cms-api/src/utils"
	"image"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"golang.org/x/exp/slices"
)

var (
	ErrReadByteFromLogoMultipartFileHeader = errors.New("read all byte from logo image multipart file header failure")
	ErrUnsupportedLogoImageType            = errors.New("resize QRCode logo to compatibility with version failure: unsupported logo image type")
	ErrContentTooLarge                     = errors.New("detect QRCode version failure: content length too large")
)

type QRCodeStorage interface {
	CreateQRCode(ctx context.Context, qrCode *entity.QRCodeCreatable) (*string, error)
}

type qrCodeBusiness struct {
	qrCodeStorage QRCodeStorage
}

func NewQRCodeBusiness(qrCodeStorage QRCodeStorage) *qrCodeBusiness {
	return &qrCodeBusiness{qrCodeStorage: qrCodeStorage}
}

func (business *qrCodeBusiness) DetectQRCodeType(content string) string {
	switch {
	case strings.Contains(content, "http") || strings.Contains(content, "https"):
		return constants.QRCodeURLType
	case strings.Contains(content, "mailto"):
		return constants.QRCodeMailType
	case strings.Contains(content, "smsto"):
		return constants.QRCodeSMSType
	case strings.Contains(content, "tel"):
		return constants.QRCodeTelType
	case strings.Contains(content, "wifi"):
		return constants.QRCodeWifiType
	default:
		return constants.QRCodeTextType
	}
}

// TODO: Take care of contentBits here:
func (business *qrCodeBusiness) DetectQRCodeVersion(content string, errorLevel int) (*int, error) {
	contentBits := len([]rune(content)) * 8
	var maximumSupportedBits []int
	if errorLevel == 1 {
		maximumSupportedBits = append(maximumSupportedBits, 152, 272, 440, 640, 864, 1088, 1248, 1856, 2192, 2592, 2960, 3424, 3688, 4184, 4712, 5176, 5768, 6360, 6888, 7456, 8048, 8752, 9392, 10208, 10960, 11744, 12248, 13048, 13880, 14744, 15640, 16568, 17528, 18448, 19472, 20528, 21616, 22496, 23648)
	} else if errorLevel == 2 {
		maximumSupportedBits = append(maximumSupportedBits, 128, 224, 352, 512, 688, 864, 992, 1232, 1456, 1728, 2032, 2320, 2672, 2920, 3320, 3624, 4056, 4504, 5016, 5352, 5712, 6256, 6880, 7312, 8000, 8496, 9024, 9544, 10136, 10984, 11640, 12328, 13048, 13800, 14496, 15312, 15936, 16816, 17728, 18672)
	} else if errorLevel == 3 {
		maximumSupportedBits = append(maximumSupportedBits, 104, 176, 272, 384, 496, 608, 704, 880, 1056, 1232, 1440, 1648, 1952, 2088, 2360, 2600, 2936, 3176, 3560, 3880, 4096, 4544, 4912, 5312, 5744, 6032, 6464, 6968, 7288, 7880, 8264, 8920, 9368, 9848, 10288, 10832, 11408, 12016, 12656, 13328)
	} else {
		maximumSupportedBits = append(maximumSupportedBits, 72, 128, 208, 288, 368, 480, 528, 688, 800, 976, 1120, 1264, 1440, 1576, 1784, 2024, 2264, 2504, 2728, 3248, 3536, 3712, 4112, 4304, 4768, 5024, 5288, 5608, 5960, 6344, 6760, 7208, 7688, 7888, 8432, 8768, 9136, 9776, 10208)
	}
	if version, err := utils.NewSliceUtil().LastIndexLessThanSliceValue(maximumSupportedBits, contentBits); err != nil {
		return nil, err
	} else {
		return version, nil
	}
}

func (business *qrCodeBusiness) ResizeLogoWithVersion(qrCode entity.QRCodeCreatable) (image.Image, error) {
	supportedLogoType := []string{"image/jpeg", "image/png", "image/jpg"}
	logoDimension := 0.2 * float64((420 + 2*(*qrCode.BorderWidth) + 80*(*qrCode.Version-1)))
	logoFile, _ := qrCode.Logo.Open()
	if logoBytes, err := io.ReadAll(logoFile); err != nil {
		fmt.Println("Error while read all byte from logo image multipart file header: " + err.Error())
		return nil, ErrReadByteFromLogoMultipartFileHeader
	} else if mimeType := http.DetectContentType(logoBytes); !slices.Contains(supportedLogoType, mimeType) {
		fmt.Println("Error while resize QRCode logo to compatibility with version: unsupported logo image type")
		return nil, ErrUnsupportedLogoImageType
	} else {
		logoImage, _, _ := image.Decode(bytes.NewReader(logoBytes))
		resizedLogo := resize.Resize(uint(logoDimension), 0, logoImage, resize.Lanczos2)
		return resizedLogo, nil
	}
}

func (business *qrCodeBusiness) Standardized(qrCode *entity.QRCodeCreatable) ([]qrcode.EncodeOption, []standard.ImageOption, error) {
	currentDir, _ := os.Getwd()
	defaultLocalFileType := os.Getenv("LOCAL_QRCODE_FILE_TYPE")
	localStorageBasePath := os.Getenv("LOCAL_QRCODE_STORAGE_DIR")
	localFileName := uuid.NewString() + "." + defaultLocalFileType
	localFilePath := filepath.Join(currentDir, localStorageBasePath, localFileName)
	qrCode.FilePath = localFilePath
	qrCode.Type = business.DetectQRCodeType(qrCode.EncodeContent)

	qrCodeConfigs := []qrcode.EncodeOption{
		qrcode.WithEncodingMode(qrcode.EncModeByte),
	}
	writerConfigs := []standard.ImageOption{
		standard.WithBorderWidth(*qrCode.BorderWidth),
	}

	if *qrCode.ErrorLevel == 1 {
		qrCodeConfigs = append(qrCodeConfigs, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionLow))
	} else if *qrCode.ErrorLevel == 2 {
		qrCodeConfigs = append(qrCodeConfigs, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionMedium))
	} else if *qrCode.ErrorLevel == 3 {
		qrCodeConfigs = append(qrCodeConfigs, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionQuart))
	} else {
		qrCodeConfigs = append(qrCodeConfigs, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest))
	}
	if qrCode.CircleShape != nil && *qrCode.CircleShape {
		writerConfigs = append(writerConfigs, standard.WithCircleShape())
	}
	if qrCode.TransparentBackground != nil && *qrCode.TransparentBackground {
		writerConfigs = append(writerConfigs, standard.WithBgTransparent())
	}

	if version, err := business.DetectQRCodeVersion(*qrCode.Content, *qrCode.ErrorLevel); err != nil {
		fmt.Println("Error while detect QR Code version: " + err.Error())
		return nil, nil, ErrContentTooLarge
	} else {
		qrCodeVersion := *version + 1
		qrCode.Version = &qrCodeVersion
		qrCodeConfigs = append(qrCodeConfigs, qrcode.WithVersion(qrCodeVersion))
	}
	if qrCode.Logo != nil {
		if logoImage, err := business.ResizeLogoWithVersion(*qrCode); err != nil {
			return nil, nil, err
		} else {
			writerConfigs = append(writerConfigs, standard.WithLogoImage(logoImage))
		}
		writerConfigs = append(writerConfigs, standard.WithBgColorRGBHex(*qrCode.Background))
		writerConfigs = append(writerConfigs, standard.WithFgColorRGBHex(*qrCode.Foreground))
	}
	return qrCodeConfigs, writerConfigs, nil
}

func (business *qrCodeBusiness) CreateQRCode(ctx context.Context, cld *cloudinary.Cloudinary, qrCode *entity.QRCodeCreatable) (*string, *string, error) {
	qrCode.Mask()
	if qrCodeConfigs, writerConfigs, err := business.Standardized(qrCode); err != nil {
		return nil, nil, err
	} else if localFilePath, err := utils.NewQrCodeUtil().SaveQRCode2LocalStorage(qrCode, qrCodeConfigs, writerConfigs); err != nil {
		return nil, nil, err
	} else if encode, err := utils.NewImageUtil().Image2Base64(*localFilePath); err != nil {
		return localFilePath, nil, err
	} else {
		qrCode.EncodeContent = *encode
		cloudinaryStorage := storage.NewCloudinaryStore(cld)
		if uploadResult, err := cloudinaryStorage.UploadSingleImage(ctx, qrCode.EncodeContent); err != nil {
			return localFilePath, nil, err
		} else {
			qrCode.PublicURL = uploadResult.URL
			if _, err := business.qrCodeStorage.CreateQRCode(ctx, qrCode); err != nil {
				return localFilePath, nil, err
			}
			return encode, &uploadResult.URL, nil
		}
	}
}

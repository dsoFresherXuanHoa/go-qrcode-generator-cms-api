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
	"reflect"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

var (
	ErrReadByteFromLogoMultipartFileHeader = errors.New("read all byte from logo image multipart file header failure")
	ErrUnsupportedLogoImageType            = errors.New("resize qrcode logo to compatibility with version failure: unsupported logo image type")
	ErrContentTooLarge                     = errors.New("detect qrcode version failure: content length too large")
	ErrQrCodeUUIDFormat                    = errors.New("validate qrcode uuid failure")
)

type QRCodeStorage interface {
	CreateQRCode(ctx context.Context, qrCode *entity.QRCodeCreatable) (*string, error)
	FindQRCodeByUUID(ctx context.Context, uuid string) (*entity.QRCodeResponse, error)
	FindQRCodeByCondition(ctx context.Context, cond map[string]interface{}, timeStat map[string]string, paging *entity.Paginate) ([]entity.QRCodeResponse, error)
}

type RedisStorage interface {
	GetQRCodeEncodeFromRedis(key string) ([]string, error)
	GetRedisKey(qrCode *entity.QRCodeCreatable) string
	SaveQRCode(qrCode *entity.QRCodeCreatable) (*string, error)
	SaveAccessToken(key string, accessToken string) error
	DeleteAccessToken(key string) error
}

type qrCodeBusiness struct {
	qrCodeStorage QRCodeStorage
	redisStorage  RedisStorage
}

func NewQRCodeBusiness(qrCodeStorage QRCodeStorage, redisStorage RedisStorage) *qrCodeBusiness {
	return &qrCodeBusiness{qrCodeStorage: qrCodeStorage, redisStorage: redisStorage}
}

func (business *qrCodeBusiness) DetectQRCodeType(content string) string {
	content = strings.ToLower(content)
	switch {
	case strings.Contains(content, "http:") || strings.Contains(content, "https"):
		return constants.QRCodeURLType
	case strings.Contains(content, "mailto:"):
		return constants.QRCodeMailType
	case strings.Contains(content, "smsto:"):
		return constants.QRCodeSMSType
	case strings.Contains(content, "tel:"):
		return constants.QRCodeTelType
	case strings.Contains(content, "wifi:"):
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
	qrCode.Type = business.DetectQRCodeType(*qrCode.Content)

	qrCodeConfigs := []qrcode.EncodeOption{
		qrcode.WithEncodingMode(qrcode.EncModeByte),
	}
	writerConfigs := []standard.ImageOption{}

	if qrCode.BorderWidth != nil {
		writerConfigs = append(writerConfigs, standard.WithBorderWidth(*qrCode.BorderWidth))
	}
	if qrCode.Background != nil {
		writerConfigs = append(writerConfigs, standard.WithBgColorRGBHex(*qrCode.Background))
	}
	if qrCode.Foreground != nil {
		writerConfigs = append(writerConfigs, standard.WithFgColorRGBHex(*qrCode.Foreground))
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
	// TODO: Transparent background not working
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
			*qrCode.Version = 10
			qrCodeConfigs = append(qrCodeConfigs, qrcode.WithVersion(*qrCode.Version))
			writerConfigs = append(writerConfigs, standard.WithLogoImage(logoImage))
		}
	} else if qrCode.Halftone != nil {
		// TODO: Halftone not working
		if localHalftonePath, err := utils.NewImageUtil().ImageMultipartFile2LocalStorage(qrCode.Halftone); err != nil {
			return nil, nil, err
		} else {
			/*
				*qrCode.Version = 20
				qrCodeConfigs = append(qrCodeConfigs, qrcode.WithVersion(*qrCode.Version))
			*/
			writerConfigs = append(writerConfigs, standard.WithHalftone(*localHalftonePath))
			writerConfigs = append(writerConfigs, standard.WithBgColorRGBHex("#000000"))
			writerConfigs = append(writerConfigs, standard.WithBgColorRGBHex("#ffffff"))
		}
	}
	return qrCodeConfigs, writerConfigs, nil
}

func (business *qrCodeBusiness) CreateQRCode(ctx context.Context, cld *cloudinary.Cloudinary, qrCode *entity.QRCodeCreatable) ([]string, []string, error) {
	var qrCodesEncode []string
	var publicUrls []string
	totalQRCode := len(qrCode.Contents)
	for i := 0; i < totalQRCode; i++ {
		qrCodeClone := *qrCode
		qrCodeClone.Mask()
		qrCodeClone.Content = &qrCode.Contents[i]
		key := business.redisStorage.GetRedisKey(&qrCodeClone)

		if redisResult, err := business.redisStorage.GetQRCodeEncodeFromRedis(key); err != storage.ErrGetQRCodeFromRedis {
			qrCodeEncode := redisResult[0]
			publicUrl := redisResult[1]

			qrCodeClone.EncodeContent = qrCodeEncode
			qrCodeClone.PublicURL = publicUrl
			qrCodeClone.Model = gorm.Model{}
			if _, err := business.qrCodeStorage.CreateQRCode(ctx, &qrCodeClone); err != nil {
				return nil, nil, err
			} else {
				qrCodesEncode = append(qrCodesEncode, qrCodeEncode)
				publicUrls = append(publicUrls, publicUrl)
			}
		} else if qrCodeConfigs, writerConfigs, err := business.Standardized(&qrCodeClone); err != nil {
			return nil, nil, err
		} else if localFilePath, err := utils.NewQrCodeUtil().SaveQRCode2LocalStorage(&qrCodeClone, qrCodeConfigs, writerConfigs); err != nil {
			return nil, nil, err
		} else if encode, err := utils.NewImageUtil().Image2Base64(*localFilePath); err != nil {
			return nil, nil, err
		} else {
			qrCodeClone.EncodeContent = *encode
			cloudinaryStorage := storage.NewCloudinaryStore(cld)
			if uploadResult, err := cloudinaryStorage.UploadSingleImage(ctx, qrCodeClone.EncodeContent); err != nil {
				return nil, nil, err
			} else {
				qrCodeClone.PublicURL = uploadResult.URL
				qrCodeClone.Content = &qrCode.Contents[i]
				if _, err := business.qrCodeStorage.CreateQRCode(ctx, &qrCodeClone); err != nil {
					return nil, nil, err
				}
				qrCodesEncode = append(qrCodesEncode, qrCodeClone.EncodeContent)
				publicUrls = append(publicUrls, qrCodeClone.PublicURL)
			}
		}
	}
	return qrCodesEncode, publicUrls, nil
}

func (business *qrCodeBusiness) FindQRCodeByUUID(ctx context.Context, qrCodeUUID string) (*entity.QRCodeResponse, error) {
	if _, err := uuid.Parse(qrCodeUUID); err != nil {
		fmt.Println("Error while validate qrCodeUUID in request: " + err.Error())
		return nil, ErrQrCodeUUIDFormat
	} else if qrCode, err := business.qrCodeStorage.FindQRCodeByUUID(ctx, qrCodeUUID); err != nil {
		fmt.Println(reflect.TypeOf(err))
		return nil, err
	} else {
		return qrCode, nil
	}
}

func (business *qrCodeBusiness) FindQRCodeByCondition(ctx context.Context, cond map[string]interface{}, timeStat map[string]string, paging *entity.Paginate) ([]entity.QRCodeResponse, error) {

	paging.Standardized()
	if qrCodes, err := business.qrCodeStorage.FindQRCodeByCondition(ctx, cond, timeStat, paging); err != nil {
		return nil, err
	} else {
		return qrCodes, nil
	}
}

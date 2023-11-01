package constants

var (
	QRCodeTextType = "text"
	QRCodeURLType  = "url"
	QRCodeMailType = "mail"
	QRCodeTelType  = "tel"
	QRCodeSMSType  = "sms"
	QRCodeWifiType = "wifi"

	InvalidQrCodeCreatableRequestFormat = "Invalid qrcode creatable request format: check your request and try again later!"

	ErrCreateQrCode = "Error while create new QrCode: make sure you has permission to do this and try again later!"

	CreateQrCodeSuccess = "Create new QrCode success: congrats!!!"
)

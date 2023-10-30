package utils

import (
	"fmt"
	"go-qrcode-generator-cms-api/src/entity"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type mailUtil struct {
	message *gomail.Message
	dialer  *gomail.Dialer
}

func NewMailUtil() *mailUtil {
	return &mailUtil{message: nil, dialer: nil}
}

func (m mailUtil) SendActivationRequestEmail(usr entity.UserCreatable) error {
	mailSenderHost := os.Getenv("MAIL_SENDER_HOST")
	mailSenderPort, _ := strconv.Atoi(os.Getenv("MAIL_SENDER_PORT"))
	mailSenderHostAddress := os.Getenv("MAIL_SENDER_HOST_ADDRESS")
	mailSenderApplicationPassword := os.Getenv("MAIL_SENDER_APPLICATION_PASSWORD")
	baseApiUrl := os.Getenv("BASE_API_URL")

	activationUrl := baseApiUrl + "/auth/activation?activationCode=" + usr.ActivationCode
	body := fmt.Sprintf(`
		Welcome: %s to QrCode Generator CMS service:
		We has been received register request from your Email address.
		If you ready send this request, please activate your account by click to below URL to use our service.
		%s
		If you not send this request, please skip this email!
		
		Thank for your caring!!!`, *usr.FirstName, activationUrl)

	m.message = gomail.NewMessage()
	m.message.SetHeader("From", mailSenderHostAddress)
	m.message.SetHeader("To", *usr.Email)
	m.message.SetHeader("Subject", "ACTIVATION REQUEST NOTIFICATION - GO QR CODE GENERATOR CMS SERVICE")
	m.message.SetBody("text/plain", body)

	m.dialer = gomail.NewDialer(mailSenderHost, mailSenderPort, mailSenderHostAddress, mailSenderApplicationPassword)
	if err := m.dialer.DialAndSend(m.message); err != nil {
		fmt.Println("Error while send activation request email to user: " + err.Error())
		return err
	}
	return nil
}

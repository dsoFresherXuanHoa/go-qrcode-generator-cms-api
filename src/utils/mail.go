package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

var (
	ErrSendActivationRequestEmail    = errors.New("send activation request email failure")
	ErrSendResetPasswordRequestEmail = errors.New("send reset password request email")
)

type mailUtil struct {
	message *gomail.Message
	dialer  *gomail.Dialer
}

func NewMailUtil() *mailUtil {
	return &mailUtil{message: nil, dialer: nil}
}

func (m mailUtil) SendActivationRequestEmail(email string, activationCode string) error {
	mailSenderHost := os.Getenv("MAIL_SENDER_HOST")
	mailSenderPort, _ := strconv.Atoi(os.Getenv("MAIL_SENDER_PORT"))
	mailSenderHostAddress := os.Getenv("MAIL_SENDER_HOST_ADDRESS")
	mailSenderApplicationPassword := os.Getenv("MAIL_SENDER_APPLICATION_PASSWORD")
	baseApiUrl := os.Getenv("BASE_API_URL")

	activationUrl := baseApiUrl + "/auth/activation?activationCode=" + activationCode
	body := fmt.Sprintf(`
		Welcome: %s
		We has been received register request from you.
		If you send this request, please activate your account by click to the below URL in order to use our service.
		%s
		In case you did not send this request, please skip this email!
		
		Thank for your caring!!!`, email, activationUrl)

	m.message = gomail.NewMessage()
	m.message.SetHeader("From", mailSenderHostAddress)
	m.message.SetHeader("To", email)
	m.message.SetHeader("Subject", "ACTIVATION REQUEST NOTIFICATION - QR CODE GENERATOR CMS")
	m.message.SetBody("text/plain", body)

	m.dialer = gomail.NewDialer(mailSenderHost, mailSenderPort, mailSenderHostAddress, mailSenderApplicationPassword)
	if err := m.dialer.DialAndSend(m.message); err != nil {
		fmt.Println("Error while send activation request email: " + err.Error())
		return ErrSendActivationRequestEmail
	}
	return nil
}

func (m mailUtil) SendResetPasswordRequestEmail(email string, activationCode string) error {
	mailSenderHost := os.Getenv("MAIL_SENDER_HOST")
	mailSenderPort, _ := strconv.Atoi(os.Getenv("MAIL_SENDER_PORT"))
	mailSenderHostAddress := os.Getenv("MAIL_SENDER_HOST_ADDRESS")
	mailSenderApplicationPassword := os.Getenv("MAIL_SENDER_APPLICATION_PASSWORD")
	baseApiUrl := os.Getenv("BASE_API_URL")

	activationUrl := baseApiUrl + "/auth/reset-password?resetCode=" + activationCode
	body := fmt.Sprintf(`
		Hi: %s
		We has been received reset password request from you.
		If you send this request, click to the below URL to reset your password:
		%s
		If you did not send this request, please skip this email!
		
		Thank for your caring!!!`, email, activationUrl)

	m.message = gomail.NewMessage()
	m.message.SetHeader("From", mailSenderHostAddress)
	m.message.SetHeader("To", email)
	m.message.SetHeader("Subject", "RESET PASSWORD REQUEST NOTIFICATION - QR CODE GENERATOR CMS")
	m.message.SetBody("text/plain", body)

	m.dialer = gomail.NewDialer(mailSenderHost, mailSenderPort, mailSenderHostAddress, mailSenderApplicationPassword)
	if err := m.dialer.DialAndSend(m.message); err != nil {
		fmt.Println("Error while send reset password request email: " + err.Error())
		return ErrSendResetPasswordRequestEmail
	}
	return nil
}

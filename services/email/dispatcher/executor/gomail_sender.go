package executor

import (
	"email/dispatcher/executor/templates"
	"github.com/go-gomail/gomail"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

var port int
var host string
var transmitter string
var pass string

type GomailSender struct{}

func init() {
	err := godotenv.Load()

	port, err = strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {

	}

	host = os.Getenv("EMAIL_HOST")
	transmitter = os.Getenv("EMAIL_SENDER")
	pass = os.Getenv("EMAIL_PASS")
}

func (s GomailSender) Send(content templates.EmailContent, email string) (err error) {
	dialer := gomail.NewDialer(host, port, transmitter, pass)
	err = dialer.DialAndSend(s.createEmail(content, email))
	return err
}

func (s GomailSender) createEmail(content templates.EmailContent, email string) *gomail.Message {
	message := gomail.NewMessage()
	message.SetHeader("From", transmitter)
	message.SetHeader("To", email)
	message.SetHeader("Subject", content.Subject)
	message.SetBody("text/plain", content.Body)

	return message
}

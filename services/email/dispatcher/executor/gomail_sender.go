package executor

import (
	"email/dispatcher/executor/templates"
	"log"
	"os"
	"strconv"

	"github.com/go-gomail/gomail"
	"github.com/joho/godotenv"
)

type GomailSender struct {
	Port        int
	Host        string
	Transmitter string
	Pass        string
}

func NewGoSender() *GomailSender {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Can not load .env config")
	}

	port, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {
		log.Fatalf("EMAIL_PORT must be integer")
	}

	return &GomailSender{
		port,
		os.Getenv("EMAIL_HOST"),
		os.Getenv("EMAIL_SENDER"),
		os.Getenv("EMAIL_PASS"),
	}
}

func (s GomailSender) Send(content templates.EmailContent, email string) (err error) {
	dialer := gomail.NewDialer(s.Host, s.Port, s.Transmitter, s.Pass)
	err = dialer.DialAndSend(s.createEmail(content, email))
	return err
}

func (s GomailSender) createEmail(content templates.EmailContent, email string) *gomail.Message {
	message := gomail.NewMessage()
	message.SetHeader("From", s.Transmitter)
	message.SetHeader("To", email)
	message.SetHeader("Subject", content.Subject)
	message.SetBody("text/plain", content.Body)

	return message
}

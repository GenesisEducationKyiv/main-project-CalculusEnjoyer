package executor

import (
	"email/config"
	"email/dispatcher/messages"

	"github.com/go-gomail/gomail"
)

type GomailSender struct {
	conf config.Config
}

func NewGoSender(conf config.Config) *GomailSender {
	return &GomailSender{
		conf,
	}
}

func (s GomailSender) Send(content messages.EmailContent, email string) (err error) {
	dialer := gomail.NewDialer(s.conf.EmailHost, s.conf.EmailPort, s.conf.EmailSender, s.conf.EmailSender)
	err = dialer.DialAndSend(s.createEmail(content, email))
	return err
}

func (s GomailSender) createEmail(content messages.EmailContent, email string) *gomail.Message {
	message := gomail.NewMessage()
	message.SetHeader("From", s.conf.EmailSender)
	message.SetHeader("To", email)
	message.SetHeader("Subject", content.Subject)
	message.SetBody("text/plain", content.Body)

	return message
}

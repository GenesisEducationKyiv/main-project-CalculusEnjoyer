package dispatcher

import (
	"email/dispatcher/messages"
)

type Sender interface {
	Send(content messages.EmailContent, email string) (err error)
}

type EmailService interface {
	SendEmail(req messages.SendEmailRequest) (err error)
}

type emailService struct{ Sender Sender }

func NewService(sender Sender) EmailService { return &emailService{sender} }

func (e emailService) SendEmail(request messages.SendEmailRequest) (err error) {
	return e.Sender.Send(request.Content, request.To)
}

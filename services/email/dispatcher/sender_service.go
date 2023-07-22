package dispatcher

import (
	domain2 "email/domain"
)

type Sender interface {
	Send(content domain2.EmailContent, email string) (err error)
}

type EmailService interface {
	SendEmail(req domain2.SendEmailRequest) (err error)
}

type emailService struct{ Sender Sender }

func NewService(sender Sender) EmailService { return &emailService{sender} }

func (e emailService) SendEmail(request domain2.SendEmailRequest) (err error) {
	return e.Sender.Send(request.Content, request.To)
}

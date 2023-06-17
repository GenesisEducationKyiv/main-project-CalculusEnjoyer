package dispatcher

import (
	sender "email/dispatcher/executor"
)

type emailService struct{ Sender sender.Sender }

func NewService(sender sender.Sender) EmailService { return &emailService{sender} }

func (e emailService) SendEmail(request sender.SendEmailRequest) (err error) {
	return e.Sender.Send(request.Content, request.To)
}

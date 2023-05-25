package dispatcher

import (
	sender "email/dispatcher/executor"
)

var dispatcher = sender.GomailSender{}

type emailService struct{}

func NewService() EmailService { return &emailService{} }

func (e emailService) SendEmail(request sender.SendEmailRequest) (err error) {
	return dispatcher.Send(request.Content, request.To)
}

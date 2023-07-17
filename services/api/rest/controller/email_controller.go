package controller

import (
	"api/domain"
	"api/rest"
	"context"
	"currency/logger"
	"net/http"
)

type EmailErrorPresenter interface {
	PresentHTTPErr(err error, w http.ResponseWriter)
}

type EmailPresenter interface {
	SuccessfulEmailsSending(w http.ResponseWriter)
	SuccessfullyAddEmail(w http.ResponseWriter)
}

type EmailService interface {
	SendRateEmails(cnx context.Context) (err error)
	AddEmail(email domain.AddEmailRequest, cnx context.Context) error
}

type EmailController struct {
	emailService EmailService
	errPresenter EmailErrorPresenter
	presenter    EmailPresenter
}

func NewEmailController(
	emailService EmailService,
	errPresenter EmailErrorPresenter,
	presenter EmailPresenter,
) *EmailController {
	return &EmailController{
		emailService: emailService,
		errPresenter: errPresenter,
		presenter:    presenter,
	}
}

func (e *EmailController) AddEmail(w http.ResponseWriter, r *http.Request) {
	logger.DefaultLog(logger.INFO, "receiving api call on add email endpoint")
	if err := r.ParseForm(); err != nil {
		e.errPresenter.PresentHTTPErr(err, w)
		return
	}

	email := r.Form.Get(rest.KeyEmail)

	if err := e.emailService.AddEmail(domain.AddEmailRequest{Email: domain.Email{Value: email}}, r.Context()); err != nil {
		logger.DefaultLog(logger.ERROR, "error while adding email")
		e.errPresenter.PresentHTTPErr(err, w)
		return
	}

	e.presenter.SuccessfullyAddEmail(w)
}

func (e *EmailController) SendBTCRateEmails(w http.ResponseWriter, r *http.Request) {
	if err := e.emailService.SendRateEmails(r.Context()); err != nil {
		e.errPresenter.PresentHTTPErr(err, w)
		return
	}

	e.presenter.SuccessfulEmailsSending(w)
}

package controller

import (
	"api/domain"
	"api/logger"
	"api/rest"
	"context"
	"net/http"
)

type EmailErrorPresenter interface {
	PresentHTTPErr(err error, w http.ResponseWriter)
}

type EmailPresenter interface {
	SuccessfulEmailsSending(w http.ResponseWriter)
	SuccessfullyAddEmail(w http.ResponseWriter)
	SuccessfullyAddEmailAndSentGreet(w http.ResponseWriter)
}

type EmailService interface {
	SendRateEmails(cnx context.Context) (err error)
	AddEmail(email domain.AddEmailRequest, cnx context.Context) error
	AddEmailWithGreeting(email domain.AddEmailRequest, ctx context.Context) error
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
	logger.DefaultLog(logger.INFO, "receiving api call on send email endpoint")
	if err := e.emailService.SendRateEmails(r.Context()); err != nil {
		logger.DefaultLog(logger.ERROR, "failed to send emails")
		e.errPresenter.PresentHTTPErr(err, w)
		return
	}

	e.presenter.SuccessfulEmailsSending(w)
}

func (e *EmailController) AddEmailWithGreetingEmail(w http.ResponseWriter, r *http.Request) {
	logger.DefaultLog(logger.INFO, "receiving api call on add email endpoint")
	if err := r.ParseForm(); err != nil {
		e.errPresenter.PresentHTTPErr(err, w)
		return
	}

	email := r.Form.Get(rest.KeyEmail)

	logger.DefaultLog(logger.INFO, "receiving api call on add email with greeting email endpoint")
	if err := e.emailService.AddEmailWithGreeting(domain.AddEmailRequest{Email: domain.Email{Value: email}}, r.Context()); err != nil {
		logger.DefaultLog(logger.ERROR, "failed to send emails")
		e.errPresenter.PresentHTTPErr(err, w)
		return
	}

	e.presenter.SuccessfullyAddEmailAndSentGreet(w)
}

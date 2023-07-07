package controller

import (
	"api/domain"
	"api/rest"
	"context"
	"net/http"
)

type EmailErrorTransformer interface {
	TransformToHTTPErr(err error, w http.ResponseWriter)
}

type EmailService interface {
	SendRateEmails(cnx context.Context) (err error)
	AddEmail(email domain.AddEmailRequest, cnx context.Context) error
}

type EmailController struct {
	emailService   EmailService
	errTransformer EmailErrorTransformer
}

func NewEmailController(
	emailService EmailService,
	errTransformer EmailErrorTransformer,
) *EmailController {
	return &EmailController{
		emailService:   emailService,
		errTransformer: errTransformer,
	}
}

func (e *EmailController) AddEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		e.errTransformer.TransformToHTTPErr(err, w)
		return
	}

	email := r.Form.Get(rest.KeyEmail)

	if err := e.emailService.AddEmail(domain.AddEmailRequest{Email: domain.Email{Value: email}}, r.Context()); err != nil {
		e.errTransformer.TransformToHTTPErr(err, w)
		return
	}
}

func (e *EmailController) SendBTCRateEmails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := e.emailService.SendRateEmails(r.Context()); err != nil {
		e.errTransformer.TransformToHTTPErr(err, w)
		return
	}
}

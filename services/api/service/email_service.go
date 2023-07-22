package service

import (
	"api/aerror"
	"api/domain"
	"api/template"
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type EmailValidator interface {
	Validate(email string) bool
}

type EmailExecutor interface {
	SendEmail(request domain.SendEmailsRequest, cnx context.Context) error
}

type EmailRepository interface {
	AddEmail(request domain.AddEmailRequest, cnx context.Context) error
	GetAllEmails(cnx context.Context) ([]domain.Email, error)
}

type EmailService struct {
	emailValidator      EmailValidator
	rateProvider        RateProvider
	emailExecutor       EmailExecutor
	storageOrchestrator EmailRepository
}

func NewEmailService(
	emailValidator EmailValidator,
	rateProvider RateProvider,
	emailExecutor EmailExecutor,
	storageOrchestrator EmailRepository,
) *EmailService {
	return &EmailService{
		emailValidator:      emailValidator,
		rateProvider:        rateProvider,
		emailExecutor:       emailExecutor,
		storageOrchestrator: storageOrchestrator,
	}
}

func (e *EmailService) AddEmail(email domain.AddEmailRequest, cnx context.Context) error {
	if !e.emailValidator.Validate(email.Email.Value) {
		return aerror.ErrInvalidEmail
	}

	if err := e.storageOrchestrator.AddEmail(email, cnx); err != nil {
		return errors.Wrap(err, "can not add email")
	}

	return nil
}

func (e *EmailService) SendRateEmails(cnx context.Context) (err error) {
	rateResp, err := e.rateProvider.GetRate(domain.RateRequest{BaseCurrency: "bitcoin", TargetCurrency: "uah"}, cnx)
	if err != nil {
		return errors.Wrap(err, "can not get rate")
	}
	rate := rateResp.Rate

	emailsResponse, err := e.storageOrchestrator.GetAllEmails(cnx)
	if err != nil {
		return errors.Wrap(err, "can not get all emails")
	}

	for i := range emailsResponse {
		if err = e.emailExecutor.SendEmail(generateRateEmail(emailsResponse[i], rate), cnx); err != nil {
			err = errors.Wrap(err, fmt.Sprintf("can not sent email to: %s", emailsResponse[i].Value))
		}
	}

	return err
}

func generateRateEmail(interceptor domain.Email, rate float64) domain.SendEmailsRequest {
	return domain.SendEmailsRequest{
		Interceptor: interceptor,
		Template: domain.EmailContent{
			Body:    template.BTCRateString + strconv.FormatFloat(rate, 'f', -1, 64),
			Subject: template.BTCRateSubject,
		},
	}
}

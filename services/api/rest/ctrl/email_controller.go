package ctrl

import (
	"api/aerror"
	"api/models"
	"api/rest"
	"api/template"
	"context"
	"email/dispatcher/messages"
	"net/http"
	"strconv"
)

type EmailValidator interface {
	Validate(email string) bool
}

type EmailExecutor interface {
	SendEmail(request models.SendEmailsRequest, cnx context.Context) error
}

type EmailRepository interface {
	AddEmail(request models.AddEmailRequest, cnx context.Context) error
	GetAllEmails(cnx context.Context) ([]models.Email, error)
}

type EmailController struct {
	emailValidator      EmailValidator
	rateProvider        CurrencyProvider
	emailExecutor       EmailExecutor
	storageOrchestrator EmailRepository
	errTransformer      ErrorTransformer
}

func NewEmailController(
	emailValidator EmailValidator,
	rateProvider CurrencyProvider,
	emailExecutor EmailExecutor,
	storageOrchestrator EmailRepository,
	errTransformer ErrorTransformer,
) *EmailController {
	return &EmailController{
		emailValidator:      emailValidator,
		rateProvider:        rateProvider,
		emailExecutor:       emailExecutor,
		storageOrchestrator: storageOrchestrator,
		errTransformer:      errTransformer,
	}
}

func (e *EmailController) AddEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		e.errTransformer.TransformToHTTPErr(err, w)
		return
	}

	email := r.Form.Get(rest.KeyEmail)
	if !e.emailValidator.Validate(email) {
		http.Error(w, aerror.ErrInvalidEmail.Error(), http.StatusBadRequest)
		return
	}

	request := models.AddEmailRequest{Email: models.Email{Value: email}}
	if err := e.storageOrchestrator.AddEmail(request, r.Context()); err != nil {
		e.errTransformer.TransformToHTTPErr(err, w)
		return
	}
}

func (e *EmailController) SendBTCRateEmails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rateResp, err := e.rateProvider.GetRate(models.RateRequest{BaseCurrency: "bitcoin", TargetCurrency: "uah"}, r.Context())
	if err != nil {
		e.errTransformer.TransformToHTTPErr(err, w)
		return
	}
	rate := rateResp.Rate

	emailsResponse, err := e.storageOrchestrator.GetAllEmails(r.Context())
	if err != nil {
		e.errTransformer.TransformToHTTPErr(err, w)
	}

	for i := range emailsResponse {
		if err = e.emailExecutor.SendEmail(generateRateEmail(emailsResponse[i], rate), r.Context()); err != nil {
			e.errTransformer.TransformToHTTPErr(err, w)
		}
	}
}

func generateRateEmail(interceptor models.Email, rate float64) models.SendEmailsRequest {
	return models.SendEmailsRequest{
		Interceptor: interceptor,
		Template: messages.EmailContent{
			Body:    template.BTCRateString + strconv.FormatFloat(rate, 'f', -1, 64),
			Subject: template.BTCRateSubject,
		},
	}
}

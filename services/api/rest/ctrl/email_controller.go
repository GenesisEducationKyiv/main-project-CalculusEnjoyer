package ctrl

import (
	"api/models"
	"api/rest"
	"api/template"
	"api/validator"
	"email/dispatcher/executor/templates"
	"net/http"
	"strconv"
)

type EmailController struct {
	emailValidator      validator.EmailValidator
	rateProvider        CurrencyProvider
	emailExecutor       EmailExecutor
	storageOrchestrator StorageOrchestrator
	errTransformer      ErrorTransformer
}

func NewEmailController(
	emailValidator validator.EmailValidator,
	rateProvider CurrencyProvider,
	emailExecutor EmailExecutor,
	storageOrchestrator StorageOrchestrator,
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

	err := r.ParseForm()
	if err != nil {
		e.errTransformer.TransformToHTTPErr(err, w)
		return
	}

	email := r.Form.Get(rest.KeyEmail)
	if !e.emailValidator.Validate(email) {
		e.errTransformer.TransformToHTTPErr(err, w)
		return
	}

	err = e.storageOrchestrator.AddEmail(models.AddEmailRequest{Email: models.Email{Value: email}}, r.Context())
	if err != nil {
		e.errTransformer.TransformToHTTPErr(err, w)
		return
	}
}

func (e *EmailController) SendEmails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rateResp, err := e.rateProvider.GetRate(models.RateRequest{BaseCurrency: "bitcoin", TargetCurrency: "uah"}, r.Context())
	if err != nil {
		e.errTransformer.TransformToHTTPErr(err, w)
		return
	}

	rate := rateResp.Rate

	emailsResponse, err := e.storageOrchestrator.GetAllEmails(r.Context())

	for i := range emailsResponse {
		err = e.emailExecutor.SendEmail(models.SendEmailsRequest{
			Interceptor: emailsResponse[i],
			Template: templates.EmailContent{
				Body:    template.BTCRateString + strconv.FormatFloat(rate, 'f', -1, 64),
				Subject: template.BTCRateSubject,
			},
		}, r.Context())
	}

	if err != nil {
		e.errTransformer.TransformToHTTPErr(err, w)
	}
}

package ctrl

import (
	"api/aerror"
	"api/models"
	"api/rest"
	"api/template"
	"api/validator"
	"context"
	"email/dispatcher/executor/templates"
	"net/http"
	"strconv"

	"github.com/fullstorydev/grpchan/httpgrpc"
	"google.golang.org/grpc/status"
)

type EmailController struct {
	emailValidator      validator.EmailValidator
	rateProvider        CurrencyProvider
	emailExecutor       EmailExecutor
	storageOrchestrator StorageOrchestrator
}

func NewEmailController(
	emailValidator validator.EmailValidator,
	rateProvider CurrencyProvider,
	emailExecutor EmailExecutor,
	storageOrchestrator StorageOrchestrator,
) *EmailController {
	return &EmailController{
		emailValidator:      emailValidator,
		rateProvider:        rateProvider,
		emailExecutor:       emailExecutor,
		storageOrchestrator: storageOrchestrator,
	}
}

func (e *EmailController) AddEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, aerror.ErrRequest.Error(), http.StatusBadRequest)
		return
	}

	email := r.Form.Get(rest.KeyEmail)
	if !e.emailValidator.Validate(email) {
		http.Error(w, aerror.ErrInvalidEmail.Error(), http.StatusBadRequest)
		return
	}

	err = e.storageOrchestrator.AddEmail(models.AddEmailRequest{Email: models.Email{Value: email}})
	if err != nil {
		httpgrpc.DefaultErrorRenderer(context.Background(), status.Convert(err), w)
		return
	}
}

func (e *EmailController) SendEmails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rateResp, err := e.rateProvider.GetRate(&models.RateRequest{BaseCurrency: "bitcoin", TargetCurrency: "uah"})
	if err != nil {
		httpgrpc.DefaultErrorRenderer(context.Background(), status.Convert(err), w)
		return
	}

	rate := rateResp.Rate

	emailsResponse, err := e.storageOrchestrator.GetAllEmails()

	for i := range emailsResponse {
		err = e.emailExecutor.SendEmail(models.SendEmailsRequest{
			Interceptor: emailsResponse[i],
			Template: templates.EmailContent{
				Body:    template.BTCRateString + strconv.FormatFloat(rate, 'f', -1, 64),
				Subject: template.BTCRateSubject,
			},
		})
	}

	if err != nil {
		httpgrpc.DefaultErrorRenderer(context.Background(), status.Convert(err), w)
	}
}

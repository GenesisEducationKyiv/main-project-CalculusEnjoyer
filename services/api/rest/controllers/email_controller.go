package controllers

import (
	APIerror "api/error"
	"api/grpc/client/currency"
	"api/grpc/client/email"
	"api/grpc/client/storage"
	"api/template"
	"api/validator"
	rateProto "currency/rate/messages/proto"
	emailProto "email/dispatcher/messages/proto"
	"net/http"
	"storage/emails/messages/proto"
	"strconv"
)

type EmailController struct {
	rateGRPCClient    currency.CurrencyGRPCClient
	emailGRPCClient   email.EmailGRPCClient
	storageGRPCClient storage.StorageGRPCClient
}

func NewEmailController() *EmailController {
	return &EmailController{
		rateGRPCClient:    *currency.NewCurrencyGRPCClient(),
		emailGRPCClient:   *email.NewEmailGRPCClient(),
		storageGRPCClient: *storage.NewStorageGRPCClient(),
	}
}

func (e *EmailController) AddEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, APIerror.InvalidRequest().Message, http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	if !validator.ValidateEmail(email) {
		http.Error(w, APIerror.InvalidEmail().Error(), http.StatusBadRequest)
		return
	}

	_, err = e.storageGRPCClient.AddEmail(&proto.AddEmailRequest{
		Email: email,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
}

func (e *EmailController) SendEmails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rateResp, err := e.rateGRPCClient.GetRate(&rateProto.RateRequest{
		BaseCurrency:   "bitcoin",
		TargetCurrency: "uah",
	})
	if err != nil {
		http.Error(w, APIerror.RateProviderError().Error(), http.StatusConflict)
		return
	}

	rate := rateResp.Rate

	emailsResponse := e.storageGRPCClient.GetAllEmails(&proto.GetAllEmailsRequest{})
	emails := emailsResponse.Email

	for i := range emails {
		err = e.emailGRPCClient.SendEmail(&emailProto.SendEmailRequest{
			Body:    template.BtcRateString + strconv.FormatFloat(rate, 'f', -1, 64),
			Subject: template.BtcRateSubject,
			To:      emails[i],
		})
	}

	if err != nil {
		http.Error(w, "some emails were not sent because of incorrect email structure", http.StatusBadRequest)
	}
}

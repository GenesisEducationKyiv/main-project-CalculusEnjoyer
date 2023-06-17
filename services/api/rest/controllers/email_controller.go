package controllers

import (
	"api/grpc/client/currency"
	"api/grpc/client/email"
	"api/grpc/client/storage"
	"api/rest/utils"
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
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	if !utils.ValidateEmail(email) {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	_, err = e.storageGRPCClient.AddEmail(proto.AddEmailRequest{
		Email: email,
	})

	if err != nil {
		http.Error(w, "email already exists in database", http.StatusConflict)
		return
	}
}

func (e *EmailController) SendEmails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rateResp, _ := e.rateGRPCClient.GetRate(rateProto.RateRequest{
		BaseCurrency:   "bitcoin",
		TargetCurrency: "uah",
	})
	rate := rateResp.Rate

	emailsResponse := e.storageGRPCClient.GetAllEmails(proto.GetAllEmailsRequest{})
	emails := emailsResponse.Email

	var err error
	for i := range emails {
		err = e.emailGRPCClient.SendEmail(emailProto.SendEmailRequest{
			Body:    utils.BtcRateString + strconv.FormatFloat(rate, 'f', -1, 64),
			Subject: utils.BtcRateSubject,
			To:      emails[i],
		})
	}

	if err != nil {
		http.Error(w, "some emails were not sent because of incorrect email structure", http.StatusBadRequest)
	}
}

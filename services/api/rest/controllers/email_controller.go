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

var rateGRPCClient currency.CurrencyGRPCClient
var emailGRPCClient email.EmailGRPCClient
var storageGRPCClient storage.StorageGRPCClient

func init() {
	rateGRPCClient = currency.CurrencyGRPCClient{}
	emailGRPCClient = email.EmailGRPCClient{}
	storageGRPCClient = storage.StorageGRPCClient{}
}

func AddEmail(w http.ResponseWriter, r *http.Request) {
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

	_, err = storageGRPCClient.AddEmail(proto.AddEmailRequest{
		Email: email,
	})

	if err != nil {
		http.Error(w, "email already exists in database", http.StatusConflict)
		return
	}

	return
}

func SendEmails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rateResp, _ := rateGRPCClient.GetRate(rateProto.RateRequest{
		BaseCurrency:   "bitcoin",
		TargetCurrency: "uah",
	})
	rate := rateResp.Rate

	emailsResponse := storageGRPCClient.GetAllEmails(proto.GetAllEmailsRequest{})
	emails := emailsResponse.Email

	var err error
	for i := range emails {
		err = emailGRPCClient.SendEmail(emailProto.SendEmailRequest{
			Body:    utils.BtcRateString + strconv.FormatFloat(rate, 'f', -1, 64),
			Subject: utils.BtcRateSubject,
			To:      emails[i],
		})
	}

	if err != nil {
		http.Error(w, "some emails were not sent because of incorrect email structure", http.StatusBadRequest)
	}

	return
}

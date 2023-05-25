package controllers

import (
	"api/grpc/client/currency"
	"currency/rate/messages/proto"
	"encoding/json"
	"net/http"
)

var currencyGRPCClient currency.CurrencyGRPCClient

func init() {
	currencyGRPCClient = currency.CurrencyGRPCClient{}
}

func GetRate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response, err := currencyGRPCClient.GetRate(proto.RateRequest{BaseCurrency: "bitcoin", TargetCurrency: "uah"})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(response.Rate)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

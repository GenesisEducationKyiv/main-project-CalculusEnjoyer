package controllers

import (
	"api/grpc/client/currency"
	"currency/rate/messages/proto"
	"encoding/json"
	"net/http"
)

type RateController struct {
	currencyGRPCClient currency.CurrencyGRPCClient
}

func NewRateController() *RateController {
	return &RateController{currencyGRPCClient: *currency.NewCurrencyGRPCClient()}
}

func (rc *RateController) GetRate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response, err := rc.currencyGRPCClient.GetRate(proto.RateRequest{BaseCurrency: "bitcoin", TargetCurrency: "uah"})

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

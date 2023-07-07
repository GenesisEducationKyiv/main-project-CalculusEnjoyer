package controller

import (
	"api/domain"
	"context"
	"encoding/json"
	"net/http"
)

type RateErrorTransformer interface {
	TransformToHTTPErr(err error, w http.ResponseWriter)
}

type RateService interface {
	GetRate(request domain.RateRequest, cnx context.Context) (*domain.RateResponse, error)
}

type RateController struct {
	rateService    RateService
	errTransformer RateErrorTransformer
}

func NewRateController(rateService RateService, errTransformer RateErrorTransformer) *RateController {
	return &RateController{rateService: rateService, errTransformer: errTransformer}
}

func (rc *RateController) GetRate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response, err := rc.rateService.GetRate(domain.RateRequest{BaseCurrency: "bitcoin", TargetCurrency: "uah"}, r.Context())
	if err != nil {
		rc.errTransformer.TransformToHTTPErr(err, w)
		return
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		rc.errTransformer.TransformToHTTPErr(err, w)
		return
	}
}

package controller

import (
	"api/domain"
	"context"
	"net/http"
)

type RateErrorPresenter interface {
	PresentHTTPErr(err error, w http.ResponseWriter)
}

type RatePresenter interface {
	SuccessfulRateResponse(w http.ResponseWriter, response domain.RateResponse)
}

type RateService interface {
	GetRate(request domain.RateRequest, cnx context.Context) (*domain.RateResponse, error)
}

type RateController struct {
	rateService  RateService
	errPresenter RateErrorPresenter
	presenter    RatePresenter
}

func NewRateController(rateService RateService, errPresenter RateErrorPresenter, presenter RatePresenter) *RateController {
	return &RateController{rateService: rateService, errPresenter: errPresenter, presenter: presenter}
}

func (rc *RateController) GetRate(w http.ResponseWriter, r *http.Request) {
	response, err := rc.rateService.GetRate(domain.RateRequest{BaseCurrency: "bitcoin", TargetCurrency: "uah"}, r.Context())
	if err != nil {
		rc.errPresenter.PresentHTTPErr(err, w)
		return
	}

	rc.presenter.SuccessfulRateResponse(w, *response)
}

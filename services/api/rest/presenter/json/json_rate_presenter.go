package json

import (
	"api/domain"
	"net/http"
)

type JSONRatePresenter struct{}

func (p *JSONRatePresenter) SuccessfulRateResponse(w http.ResponseWriter, response domain.RateResponse) {
	DecodeJSONResponse(w, response)
}

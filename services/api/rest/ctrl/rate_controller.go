package ctrl

import (
	"api/models"
	"context"
	"encoding/json"
	"net/http"
)

type ErrorTransformer interface {
	TransformToHTTPErr(err error, w http.ResponseWriter)
}

type CurrencyProvider interface {
	GetRate(request models.RateRequest, cnx context.Context) (*models.RateResponse, error)
}

type RateController struct {
	currencyGRPCClient CurrencyProvider
	errTransformer     ErrorTransformer
}

func NewRateController(provider CurrencyProvider, errTransformer ErrorTransformer) *RateController {
	return &RateController{currencyGRPCClient: provider, errTransformer: errTransformer}
}

func (rc *RateController) GetRate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response, err := rc.currencyGRPCClient.GetRate(models.RateRequest{BaseCurrency: "bitcoin", TargetCurrency: "uah"}, r.Context())
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

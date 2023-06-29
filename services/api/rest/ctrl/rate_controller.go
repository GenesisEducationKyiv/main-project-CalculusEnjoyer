package ctrl

import (
	"api/aerror"
	"api/models"
	"encoding/json"
	"net/http"

	"github.com/fullstorydev/grpchan/httpgrpc"
	"google.golang.org/grpc/status"
)

type RateController struct {
	currencyGRPCClient CurrencyProvider
}

func NewRateController(provider CurrencyProvider) *RateController {
	return &RateController{currencyGRPCClient: provider}
}

func (rc *RateController) GetRate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response, err := rc.currencyGRPCClient.GetRate(&models.RateRequest{BaseCurrency: "bitcoin", TargetCurrency: "uah"}, r.Context())
	if err != nil {
		httpgrpc.DefaultErrorRenderer(r.Context(), status.Convert(err), w)
		return
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, aerror.ErrFailedToEncodeResponse.Error(), http.StatusInternalServerError)
		return
	}
}

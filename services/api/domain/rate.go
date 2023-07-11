package domain

type RateRequest struct {
	BaseCurrency   Currency
	TargetCurrency Currency
}

type RateResponse struct {
	Rate float64 `json:"rate"`
}

package models

type RateRequest struct {
	BaseCurrency   string
	TargetCurrency string
}

type RateResponse struct {
	Rate float64 `json:"rate"`
}

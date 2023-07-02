package rate

import (
	"currency/cerror"
	"currency/rate/messages"
	"log"
	"testing"
	"time"
)

const (
	testRate          = 1000
	testUNIXTimeStamp = 1687551173
)

func TestGetRate(t *testing.T) {
	srv := NewRateService(&stubRateProvider{}, &stubTimeProvider{})

	rate, _ := srv.GetRate(messages.RateRequest{TargetCurrency: "bitcoin", BaseCurrency: "uah"})

	if rate.Rate != testRate {
		log.Fatalf(`%s: %f != %d`, "wrong result", rate.Rate, testRate)
	}
	if rate.Timestamp != time.Unix(testUNIXTimeStamp, 0) {
		log.Fatalf(`%s: %s != %s`, "wrong result", rate.Timestamp.String(), time.Unix(testUNIXTimeStamp, 0).String())
	}
}

func TestGetErrRate(t *testing.T) {
	srv := NewRateService(&stubErrorRateProvider{}, &stubTimeProvider{})

	rate, err := srv.GetRate(messages.RateRequest{TargetCurrency: "bitcoin", BaseCurrency: "uah"})

	if err == nil {
		log.Fatalf(`%s: %d`, "error is nil while it must not", err)
	}
	if rate.Rate != cerror.ErrRateValue {
		log.Fatalf(`%s: %f != %d`, "wrong result", rate.Rate, testRate)
	}
	if rate.Timestamp != time.Unix(testUNIXTimeStamp, 0) {
		log.Fatalf(`%s: %s != %s`, "wrong result", rate.Timestamp.String(), time.Unix(testUNIXTimeStamp, 0).String())
	}
}

type stubRateProvider struct{}

func (r *stubRateProvider) GetExchangeRate(baseCurrency, targetCurrency messages.Currency) (rate float64, err error) {
	return testRate, nil
}

type stubErrorRateProvider struct{}

func (r *stubErrorRateProvider) GetExchangeRate(baseCurrency, targetCurrency messages.Currency) (rate float64, err error) {
	return cerror.ErrRateValue, cerror.ErrRate
}

type stubTimeProvider struct{}

func (t *stubTimeProvider) Now() time.Time {
	return time.Unix(testUNIXTimeStamp, 0)
}

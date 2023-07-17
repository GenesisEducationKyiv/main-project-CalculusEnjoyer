package main

import (
	"api/config"
	"api/grpc/client/currency"
	"api/grpc/client/email"
	"api/grpc/client/storage"
	"api/logger"
	"api/rest"
	"api/rest/controller"
	"api/rest/presenter/json"
	"api/service"
	"api/time"
	"api/validator"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	run()
}

func run() {
	conf := config.LoadFromENV()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	loggerInstance := logger.NewBrokerLogger(&time.SystemTime{}, conf)
	loggerInstance.Init()
	logger.SetDefaultLogger(loggerInstance)

	emailService := service.NewEmailService(
		validator.NewRegexValidator(*validator.DefaultEmailRegex),
		currency.NewCurrencyGRPCClient(conf),
		email.NewEmailGRPCClient(conf),
		storage.NewStorageGRPCClient(conf))

	email := controller.NewEmailController(
		emailService,
		&json.GRPCErrHTTPPresenter{},
		&json.JSONEmailPresenter{})
	rate := controller.NewRateController(
		service.NewRateServiece(currency.NewCurrencyGRPCClient(conf)),
		&json.GRPCErrHTTPPresenter{},
		&json.JSONRatePresenter{})

	r.Route(rest.Api, func(r chi.Router) {
		r.Get(rest.BTCRate, rate.GetBTCRate)
		r.Post(rest.Rate, rate.GetRate)
		r.Post(rest.AddEmails, email.AddEmail)
		r.Post(rest.SendEmails, email.SendBTCRateEmails)
	})

	http.ListenAndServe(":"+strconv.Itoa(conf.Port), r)
}

package main

import (
	"api/config"
	"api/grpc/client"
	"api/grpc/client/currency"
	"api/grpc/client/email"
	"api/grpc/client/storage"
	"api/rest"
	"api/rest/controller"
	"api/service"
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

	emailService := service.NewEmailService(
		validator.NewRegexValidator(*validator.DefaultEmailRegex),
		currency.NewCurrencyGRPCClient(conf),
		email.NewEmailGRPCClient(conf),
		storage.NewStorageGRPCClient(conf))

	email := controller.NewEmailController(
		emailService,
		&client.GRPCErrHTTPTransformer{})
	rate := controller.NewRateController(service.NewRateServiece(currency.NewCurrencyGRPCClient(conf)), &client.GRPCErrHTTPTransformer{})

	r.Route(rest.Api, func(r chi.Router) {
		r.Get(rest.Rate, rate.GetRate)
		r.Post(rest.AddEmails, email.AddEmail)
		r.Post(rest.SendEmails, email.SendBTCRateEmails)
	})

	http.ListenAndServe(":"+strconv.Itoa(conf.Port), r)
}

package main

import (
	"api/rest"
	"api/rest/controllers"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	run()
}

func run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting the .env config")
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))

	if err != nil {
		log.Fatalf("Port must be integer")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	email := controllers.NewEmailController()
	rate := controllers.NewRateController()

	r.Get(rest.Api+rest.Rate, rate.GetRate)
	r.Post(rest.Api+rest.AddEmails, email.AddEmail)
	r.Post(rest.Api+rest.SendEmails, email.SendEmails)

	http.ListenAndServe(":"+strconv.Itoa(port), r)
}

package main

import (
	"api/rest"
	"api/rest/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
)

var port int

func main() {
	run()
}

func init() {
	godotenv.Load()
	p, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Wrong Port")
	}

	port = p
}

func run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get(rest.Api+rest.Rate, controllers.GetRate)
	r.Post(rest.Api+rest.AddEmails, controllers.AddEmail)
	r.Post(rest.Api+rest.SendEmails, controllers.SendEmails)

	http.ListenAndServe(":"+strconv.Itoa(port), r)
}

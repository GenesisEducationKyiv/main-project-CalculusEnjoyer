package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Config struct {
	Port        int
	Network     string
	EmailPass   string
	EmailSender string
	EmailHost   string
	EmailPort   int
}

func LoadFromENV() Config {
	conf := Config{}
	err := godotenv.Load(".env")
	if err != nil {
		panic(errors.Wrap(err, "Can not load config"))
	}

	conf.Network = os.Getenv("NETWORK")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(errors.Wrap(err, "Can not load PORT"))
	}
	conf.Port = port

	emailPort, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {
		panic(errors.Wrap(err, "Can not load EMAIL_PORT"))
	}
	conf.EmailPort = emailPort

	conf.EmailSender = os.Getenv("EMAIL_SENDER")
	conf.EmailHost = os.Getenv("EMAIL_HOST")
	conf.EmailPass = os.Getenv("EMAIL_PASS")

	return conf
}

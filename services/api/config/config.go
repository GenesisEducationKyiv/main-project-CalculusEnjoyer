package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Config struct {
	Port            int
	EmailNetwork    string
	EmailPort       int
	CurrencyNetwork string
	CurrencyPort    int
	StorageNetwork  string
	StoragePort     int
	KafkaAddress    string
}

func LoadFromENV() Config {
	conf := Config{}
	err := godotenv.Load(".env")
	if err != nil {
		panic(errors.Wrap(err, "Can not load config"))
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(errors.Wrap(err, "Can not load PORT"))
	}
	conf.Port = port

	email, err := strconv.Atoi(os.Getenv("EMAIL_SERVICE_PORT"))
	if err != nil {
		panic(errors.Wrap(err, "Can not load EMAIL_SERVICE_PORT"))
	}
	conf.EmailPort = email

	currency, err := strconv.Atoi(os.Getenv("CURRENCY_SERVICE_PORT"))
	if err != nil {
		panic(errors.Wrap(err, "Can not load CURRENCY_SERVICE_PORT"))
	}
	conf.CurrencyPort = currency

	storage, err := strconv.Atoi(os.Getenv("STORAGE_SERVICE_PORT"))
	if err != nil {
		panic(errors.Wrap(err, "Can not load STORAGE_SERVICE_PORT"))
	}
	conf.StoragePort = storage

	conf.EmailNetwork = os.Getenv("EMAIL_NETWORK")
	conf.CurrencyNetwork = os.Getenv("CURRENCY_NETWORK")
	conf.StorageNetwork = os.Getenv("STORAGE_NETWORK")
	conf.KafkaAddress = os.Getenv("KAFKA_ADDRESS")

	return conf
}

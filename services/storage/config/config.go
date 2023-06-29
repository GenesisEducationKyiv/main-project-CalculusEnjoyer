package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Config struct {
	Port             int
	EmailStoragePath string
	Network          string
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

	conf.EmailStoragePath = os.Getenv("STORAGE_PATH")
	conf.Network = os.Getenv("NETWORK")

	return conf
}

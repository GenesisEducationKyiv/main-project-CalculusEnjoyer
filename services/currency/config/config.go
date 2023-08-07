package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Config struct {
	Port           int
	Network        string
	CoinGekcoURL   string
	CoinApiURL     string
	CoinApiKey     string
	KunaURL        string
	AmqpURL        string
	CacheValidTime time.Duration
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

	conf.AmqpURL = os.Getenv("AMQP_URL")

	cachedTime, err := strconv.ParseInt(os.Getenv("RATE_CACHE_TIME"), 10, 64)
	if err != nil {
		panic(errors.Wrap(err, "Can not load PORT"))
	}
	conf.CacheValidTime = time.Duration(cachedTime * int64(time.Millisecond))

	conf.CoinGekcoURL = os.Getenv("COINGEKCO_URL")
	conf.Network = os.Getenv("NETWORK")
	conf.CoinApiKey = os.Getenv("COIN_API_KEY")
	conf.CoinApiURL = os.Getenv("COINAPI_URL")
	conf.KunaURL = os.Getenv("KUNA_URL")

	return conf
}

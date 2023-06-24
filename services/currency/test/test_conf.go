package test

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Config struct {
	Port              int
	Network           string
	TestRate          float64
	TestUnixTimeStamp time.Time
}

func LoadFromENV() Config {
	conf := Config{}
	err := godotenv.Load(".env.test")
	if err != nil {
		panic(errors.Wrap(err, "Can not load config"))
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(errors.Wrap(err, "Can not load PORT"))
	}
	conf.Port = port

	testRate, err := strconv.ParseFloat(os.Getenv("TEST_RATE"), 64)
	if err != nil {
		panic(errors.Wrap(err, "Can not load PORT"))
	}
	conf.TestRate = testRate

	testTimeStamp, err := strconv.ParseInt(os.Getenv("TEST_RATE"), 10, 64)
	if err != nil {
		panic(errors.Wrap(err, "Can not load PORT"))
	}
	conf.TestUnixTimeStamp = time.Unix(testTimeStamp, 0)

	conf.Network = os.Getenv("NETWORK")

	return conf
}

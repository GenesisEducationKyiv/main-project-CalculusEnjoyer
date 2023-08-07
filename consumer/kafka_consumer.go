package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
)

func main() {
	topic := "log"
	partition := 0

	err := godotenv.Load(".env")
	if err != nil {
		panic(errors.Wrap(err, "Can not load config"))
	}
	kafkaAddress := os.Getenv("KAFKA_ADDRESS")

	level := "DEBUG"
	if len(os.Args) > 1 {
		level = os.Args[1]
	}

	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaAddress, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	for {
		m, err := conn.ReadMessage(1e6)
		if err != nil {
			log.Printf(errors.Wrap(err, "err while reading batch").Error())
		}
		if string(m.Key) == level {
			fmt.Println(string(m.Value))
		}
	}
}

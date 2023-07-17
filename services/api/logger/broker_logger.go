package logger

import (
	"api/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type TimeProvider interface {
	Now() time.Time
}

const topicName = "log"

type BrokerLogger struct {
	timeProvider TimeProvider
	brokerUrl    string
	conn         *kafka.Conn
}

func NewBrokerLogger(timeProvider TimeProvider, conf config.Config) *BrokerLogger {
	return &BrokerLogger{
		timeProvider: timeProvider,
		brokerUrl:    conf.KafkaAddress,
	}
}

func (b *BrokerLogger) Init() {
	_, err := b.getConnection()
	if err != nil {
		log.Fatalf(errors.Wrap(err, "can not connect to the broker").Error())
	}
}

func (b *BrokerLogger) Log(level LogLevel, message string) {
	log.Printf(b.createLogMessage(level, message))
	err := b.publish(level, message)
	if err != nil {
		log.Printf(err.Error())
	}
}

func (b *BrokerLogger) getConnection() (*kafka.Conn, error) {
	if b.conn != nil && b.checkForHealth() {
		return b.conn, nil
	}

	conn, err := kafka.DialLeader(context.Background(), "tcp", b.brokerUrl, topicName, 0)
	if err != nil {
		log.Println("can not connect to the broker")
		return nil, errors.Wrap(err, "can not connect to the broker")
	}

	b.conn = conn
	return conn, nil
}

func (b *BrokerLogger) checkForHealth() bool {
	if _, err := b.conn.Brokers(); err != nil {
		return false
	}
	return true
}

func (b *BrokerLogger) publish(level LogLevel, message string) error {
	conn, err := b.getConnection()
	if err != nil {
		return errors.Wrap(err, "can not access the connection")
	}

	_, err = conn.WriteMessages(b.createKafkaMessage(level, message))
	return errors.Wrap(err, "can not exude message to Kafka")
}

func (b *BrokerLogger) createLogMessage(level LogLevel, message string) string {
	currentTime := b.timeProvider.Now().Format(time.UnixDate)
	return fmt.Sprintf("[%s]: %s: %s", string(level), currentTime, message)
}

func (b *BrokerLogger) createKafkaMessage(level LogLevel, message string) kafka.Message {
	return kafka.Message{Value: []byte(b.createLogMessage(level, message)), Key: []byte(level)}
}

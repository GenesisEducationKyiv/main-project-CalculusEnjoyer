package logger

import (
	"context"
	"currency/config"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type TimeProvider interface {
	Now() time.Time
}

const defaultLogExchangeName = "log"

type BrokerLogger struct {
	timeProvider TimeProvider
	brokerUrl    string
	channel      *amqp091.Channel
	conn         *amqp091.Connection
}

func NewBrokerLogger(timeProvider TimeProvider, conf config.Config) *BrokerLogger {
	return &BrokerLogger{
		timeProvider: timeProvider,
		brokerUrl:    conf.AmqpURL,
	}
}

func (l *BrokerLogger) InitQueue() {
	conn, err := amqp091.Dial(l.brokerUrl)
	if err != nil {
		log.Println(errors.Wrap(err, "can not connect to the broker"))
	}
	l.conn = conn

	l.channel, err = l.conn.Channel()
	if err != nil {
		log.Println(errors.Wrap(err, "can not create channel to the broker"))
	}

	go listenToClosingConn(l.conn)
	go listenToClosingChan(l.channel)
}

func (l *BrokerLogger) Log(level LogLevel, message string) {
	log.Printf(l.createLogMessage(level, message))
	err := l.publish(level, message)
	if err != nil {
		log.Printf(err.Error())
	}
}

func (l *BrokerLogger) publish(level LogLevel, message string) error {
	return l.channel.PublishWithContext(
		context.Background(),
		defaultLogExchangeName,
		string(level),
		false,
		false,
		amqp091.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(l.createLogMessage(level, message)),
			DeliveryMode: amqp.Persistent,
		})
}

func listenToClosingConn(conn *amqp091.Connection) {
	notifyConnClose := conn.NotifyClose(make(chan *amqp091.Error))
	err, ok := <-notifyConnClose

	if !ok {
		notifyConnClose = nil
	} else {
		log.Printf("connection closed, error %s", err)
	}
}

func listenToClosingChan(ch *amqp091.Channel) {
	notifyChanClose := ch.NotifyClose(make(chan *amqp091.Error))
	err, ok := <-notifyChanClose

	if !ok {
		notifyChanClose = nil
	} else {
		log.Printf("chan closed, error %s", err)
	}
}

func (l *BrokerLogger) createLogMessage(level LogLevel, message string) string {
	currentTime := l.timeProvider.Now().Format(time.UnixDate)
	return fmt.Sprintf("[%s]: %s: %s", string(level), currentTime, message)
}

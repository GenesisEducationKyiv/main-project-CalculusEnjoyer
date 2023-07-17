package logger

import (
	"context"
	"currency/config"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
	"github.com/streadway/amqp"
)

type TimeProvider interface {
	Now() time.Time
}

const defaultLogExchangeName = "log"

type RabbitLogger struct {
	timeProvider TimeProvider
	brokerUrl    string
	channel      *amqp091.Channel
}

func NewBrokerLogger(timeProvider TimeProvider, conf config.Config) *RabbitLogger {
	return &RabbitLogger{
		timeProvider: timeProvider,
		brokerUrl:    conf.AmqpURL,
	}
}

func (l *RabbitLogger) Init() {
	conn, err := amqp091.Dial(l.brokerUrl)
	if err != nil {
		log.Fatalf(errors.Wrap(err, "can not connect to the broker").Error())
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf(errors.Wrap(err, "can not create channel to the broker").Error())
	}

	l.channel = channel

	go l.listenToClosingChan(l.channel)
}

func (l *RabbitLogger) Log(level LogLevel, message string) {
	log.Printf(l.createLogMessage(level, message))
	err := l.publish(level, message)

	if err != nil {
		log.Printf(err.Error())
	}
}

func (l *RabbitLogger) getChannel() (*amqp091.Channel, error) {
	if !l.channel.IsClosed() {
		return l.channel, nil
	}

	conn, err := amqp091.Dial(l.brokerUrl)
	if err != nil {
		log.Println(errors.Wrap(err, "can not connect to the broker"))
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Println(errors.Wrap(err, "can not create channel to the broker"))
		return nil, err
	}

	l.channel = channel
	return l.channel, nil
}

func (l *RabbitLogger) publish(level LogLevel, message string) error {
	channel, err := l.getChannel()
	if err != nil {
		log.Println(errors.Wrap(err, "can not load channel"))
		return errors.Wrap(err, "can not load channel")
	}

	err = channel.PublishWithContext(
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

	if err != nil {
		log.Println(errors.Wrap(err, "can not publish log"))
	}

	return nil
}

func (l *RabbitLogger) listenToClosingChan(ch *amqp091.Channel) {
	notifyChanClose := ch.NotifyClose(make(chan *amqp091.Error))
	err, ok := <-notifyChanClose

	if !ok {
		close(notifyChanClose)
		log.Printf(l.createLogMessage(INFO, "channel is closed"))
	} else {
		log.Printf(l.createLogMessage(ERROR, fmt.Sprintf("chan closed, error %s", err)))
	}
}

func (l *RabbitLogger) createLogMessage(level LogLevel, message string) string {
	currentTime := l.timeProvider.Now().Format(time.UnixDate)
	return fmt.Sprintf("[%s]: %s: %s", string(level), currentTime, message)
}

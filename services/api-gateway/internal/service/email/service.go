package email

import (
	"context"
	"encoding/json"
	"gotify/services/api-gateway/internal/config"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Context = context.Context

type emailService struct {
	log        *log.Logger
	connection *amqp.Connection
	channel    *amqp.Channel
	config     config.Config
	url        string
}

func New(l *log.Logger) *emailService {
	config := config.New()

	return &emailService{
		log:    l,
		config: config,
		url:    config.EmailServiceUrl,
	}
}

func (s *emailService) SendNotification(ctx Context, body any) {
	if s.connection == nil {
		s.connect()
	}

	bytes, err := json.Marshal(body)

	if err != nil {
		s.log.Print(err)
		return
	}

	err = s.channel.PublishWithContext(
		ctx,
		"email",
		"",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         bytes,
		},
	)

	if err != nil {
		s.log.Print(err)
		return
	}

	return
}

func (s *emailService) Close() {
	if s.connection != nil {
		s.connection.Close()
	}

	if s.channel != nil {
		s.channel.Close()
	}
}

func (s *emailService) connect() error {
	conn, err := amqp.Dial(s.url)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(
		"email",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	s.connection = conn
	s.channel = ch

	return nil
}

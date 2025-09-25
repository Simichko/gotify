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
	queue      amqp.Queue
	config     config.Config
	url        string
}

func New(l *log.Logger) *emailService {
	// defaultUrl := "amqp://guest:guest@rabbitmq:5672/"
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
		"",
		s.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
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
	// failOnError(err, "Failed to connect to RabbitMQ")
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	// failOnError(err, "Failed to open channel")
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"email",
		false,
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
	s.queue = q

	// failOnError(err, "Failed to declare a queue")

	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

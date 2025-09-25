package main

import (
	"fmt"
	"gotify/shared/app"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func sendEmail(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Send email")
	return
}

func main() {
	app := app.New()

	app.Configure(func(s *http.Server) {
		s.Addr = ":8081"
	})

	app.HandleFunc("POST /", sendEmail)

	conn, err := amqp.Dial("amqp://rmuser:rmpassword@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"email",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"email-consumer",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)
		}
	}()

	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	app.Run()
}

package main

import (
	"fmt"
	"gotify/shared/app"
	"log"
	"net/http"
	"time"

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

	err = ch.ExchangeDeclare(
		"email",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"emails",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(1, 0, false)
	failOnError(err, "Failed to set QoS")

	err = ch.QueueBind(
		q.Name,
		"",
		"email",
		false,
		nil,
	)
	failOnError(err, "Failed to bind with exchange")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)
			time.Sleep(3 * time.Second)
			log.Printf("Message sent")
			msg.Ack(false)
		}
	}()

	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	app.Run()
}

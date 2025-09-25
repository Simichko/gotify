package main

import (
	"gotify/services/api-gateway/internal/service"
	"gotify/services/api-gateway/internal/service/email"
	"gotify/shared/app"
	"log"
	"net/http"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	log := log.Default()

	app := app.New()

	app.Configure(func(s *http.Server) {
		s.Addr = ":8080"
	})

	email := email.New(log)
	defer email.Close()

	services := service.Services{
		"email": email,
	}

	apiGateway := service.New(log, services)

	app.HandleFunc("POST /", apiGateway.SendNotifications)

	app.Run()
}

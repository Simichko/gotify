package main

import (
	"gotify/services/api-gateway/internal/service"
	"gotify/shared/app"
	"log"
	"net/http"
)

func main() {
	app := app.New()

	app.Configure(func(s *http.Server) {
		s.Addr = ":8080"
	})

	apiGateway := service.New(log.Default())
	app.HandleFunc("POST /", apiGateway.SendNotifications)

	app.Run()
}

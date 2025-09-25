package service

import (
	"context"
	"encoding/json"
	"gotify/services/api-gateway/internal/config"
	"gotify/services/api-gateway/internal/types"
	"log"
	"net/http"
)

type Context = context.Context

type NotificationService interface {
	SendNotification(ctx Context, body any)
}

type Services = map[string]NotificationService

type service struct {
	log      *log.Logger
	config   config.Config
	services Services
}

func New(l *log.Logger, services Services) *service {
	return &service{
		l,
		config.New(),
		services,
	}
}

func (s *service) SendNotifications(w http.ResponseWriter, req *http.Request) {
	log := s.log
	config := s.config

	var notifications types.Notifications

	err := json.NewDecoder(req.Body).Decode(&notifications)

	if err != nil {
		log.Println("Something strange happened")
		return
	}

	if notifications.Email != nil {
		url := config.EmailServiceUrl

		log.Printf("Email sended to %v\n", url)

		body := notifications.Email

		if email, ok := s.services["email"]; ok {
			email.SendNotification(req.Context(), body)
		}

		return
	}
}

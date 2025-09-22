package service

import (
	"encoding/json"
	"gotify/services/api-gateway/internal/config"
	"gotify/services/api-gateway/internal/types"
	"io"
	"log"
	"net/http"
)

type service struct {
	log    *log.Logger
	config *config.Config
}

func New(l *log.Logger) *service {
	return &service{l, config.New()}
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

		contentType := "application/json"

		body, reqWriter := io.Pipe()

		go func() {
			enc := json.NewEncoder(reqWriter)
			data := notifications.Email

			reqWriter.CloseWithError(enc.Encode(data))
		}()

		defer body.Close()

		resp, err := http.Post(url, contentType, body)

		if err != nil {
			log.Println(err)
			return
		}

		defer resp.Body.Close()

		return
	}
}

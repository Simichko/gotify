package main

import (
	"fmt"
	"gotify/shared/app"
	"net/http"
)

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

	app.Run()
}

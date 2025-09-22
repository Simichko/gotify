package app

import (
	"context"
	"errors"
	"fmt"
	"gotify/shared/env"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type app struct {
	log    *log.Logger
	mux    *http.ServeMux
	server http.Server
}

func New() *app {
	mux := http.NewServeMux()

	defaultPort := 8000
	port := env.GetEnvAsInt("PORT", defaultPort)
	addr := fmt.Sprintf(":%d", port)

	return &app{
		log: log.Default(),
		mux: mux,
		server: http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

func (app *app) Configure(options func(*http.Server)) {
	options(&app.server)
}

func (app *app) AddLogger(l *log.Logger) {
	app.log = l
}

func (app *app) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	app.mux.HandleFunc(pattern, handler)
}

func (app *app) Run() {
	log := app.log
	server := &app.server
	bckgCtx := context.Background()

	ctx, stop := signal.NotifyContext(
		bckgCtx,
		os.Interrupt,
		syscall.SIGTERM)

	defer stop()

	go func() {
		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}

		log.Println("Stopped serving new connections.")
	}()

	log.Printf("Server started at http://*%s\n", server.Addr)
	<-ctx.Done()

	log.Println("Shutting down gracefully...")

	timeout := 10 * time.Second
	shutdownCtx, cancel := context.WithTimeout(bckgCtx, timeout)
	defer cancel()

	err := server.Shutdown(shutdownCtx)

	if err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped.")
}

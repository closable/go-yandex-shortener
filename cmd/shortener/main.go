package main

import (
	"fmt"
	"net/http"

	"github.com/closable/go-yandex-shortener/internal/config"
	"github.com/closable/go-yandex-shortener/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.ParseConfigEnv()
	config.ParseFlags()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	srvAdr := config.GetEnvParam("RUN_SERVER")

	fmt.Println("Running server on", srvAdr)
	r := chi.NewRouter()
	r.Get("/{id}", handlers.GetEndpointByShortener)
	r.Post("/", handlers.GenerateShortener)

	// return http.ListenAndServe(config.FlagRunAddr, r)
	return http.ListenAndServe(srvAdr, r)
	// return http.ListenAndServe(`:8080`, http.HandlerFunc(handlers.Webhook))
}

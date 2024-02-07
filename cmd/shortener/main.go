package main

import (
	"fmt"
	"net/http"

	"github.com/closable/go-yandex-shortener/internal/config"
	"github.com/closable/go-yandex-shortener/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.ParseFlags()
	//fmt.Println("Running server on", flagRunAddr)
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	fmt.Println("Running server on", config.FlagRunAddr)
	r := chi.NewRouter()
	r.Get("/{id}", handlers.GetEndpointByShortener)
	r.Post("/", handlers.GenerateShortener)
	return http.ListenAndServe(config.FlagRunAddr, r)

	// return http.ListenAndServe(`:8080`, http.HandlerFunc(handlers.Webhook))
}

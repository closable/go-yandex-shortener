package main

import (
	"net/http"

	"github.com/closable/go-yandex-shortener/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	r := chi.NewRouter()
	r.Get("/{id}", handlers.GetEndpointByShortener)
	r.Post("/", handlers.GenerateShortener)
	return http.ListenAndServe(`:8080`, r)

	// return http.ListenAndServe(`:8080`, http.HandlerFunc(handlers.Webhook))
}

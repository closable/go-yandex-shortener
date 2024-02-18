package main

import (
	"fmt"
	"net/http"

	"github.com/closable/go-yandex-shortener/internal/config"
	"github.com/closable/go-yandex-shortener/internal/handlers"
	"github.com/closable/go-yandex-shortener/internal/storage"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	cfg := config.LoadConfig()
	store := storage.New()
	handler := handlers.New(store, cfg.BaseURL)

	fmt.Println("Running server on", cfg.ServerAddress)

	return http.ListenAndServe(cfg.ServerAddress, handler.InitRouter())
}

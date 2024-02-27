package main

import (
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
	logger := handlers.NewLogger()
	sugar := *logger.Sugar()

	producer, err := handlers.NewProducer(cfg.FileStore)
	if err != nil {
		sugar.Infoln("An error occurred while create the file ", cfg.FileStore)
		panic(err)
	}
	consumer, err := handlers.NewConsumer(cfg.FileStore)
	if err != nil {
		sugar.Infoln("An error occurred while reading the file ", cfg.FileStore)
		panic(err)
	}
	defer producer.Close()

	handler := handlers.New(store, cfg.BaseURL, logger, producer, consumer, 1)

	sugar.Infoln("Running server on", cfg.ServerAddress)
	return http.ListenAndServe(cfg.ServerAddress, handler.InitRouter())
}

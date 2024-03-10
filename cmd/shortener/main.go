package main

import (
	"net/http"
	"os"
	"path/filepath"

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
	if len(cfg.FileStore) > 0 {
		os.MkdirAll(filepath.Dir(cfg.FileStore), os.ModePerm)
	}
	store := storage.New()
	logger := handlers.NewLogger()
	sugar := *logger.Sugar()

	sugar.Infoln("DSN configure ", cfg.DSN)

	dbms, err := storage.NewDBMS(cfg.DSN)
	if err != nil {
		sugar.Panicln("Unable to connection to database", err)
		os.Exit(1)
	}

	handler := handlers.New(store, cfg.BaseURL, logger, cfg.FileStore, dbms, 1)
	sugar.Infoln("File store path", cfg.FileStore)
	sugar.Infoln("Running server on", cfg.ServerAddress)
	return http.ListenAndServe(cfg.ServerAddress, handler.InitRouter())
}

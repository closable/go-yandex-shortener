package main

import (
	"fmt"
	"net/http"
	"os"

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
	// generate test body
	// utils.GenerateBatchBody(100000)

	cfg := config.LoadConfig()
	logger := handlers.NewLogger()
	sugar := *logger.Sugar()

	var store handlers.Storager
	var err error

	var storeMsg string
	if len(cfg.DSN) > 0 {
		store, err = storage.NewDBMS(cfg.DSN)
		storeMsg = fmt.Sprintf("Store DBMS setup successfuly -> %s", cfg.DSN)
	} else if len(cfg.FileStore) > 0 {
		store, err = storage.NewFile(cfg.FileStore)
		storeMsg = fmt.Sprintf("Store File setup successfuly -> %s", cfg.FileStore)
	} else {
		store, err = storage.NewMemory()
		storeMsg = fmt.Sprintf("Store Memory setup successfuly -> %s", "default")
	}

	if err != nil {
		sugar.Panicln("Store invalid")
		os.Exit(1)
	}

	handler := handlers.New(store, cfg.BaseURL, logger, 1)

	sugar.Infoln(storeMsg)
	sugar.Infoln("Running server on", cfg.ServerAddress)
	return http.ListenAndServe(cfg.ServerAddress, handler.InitRouter())
}

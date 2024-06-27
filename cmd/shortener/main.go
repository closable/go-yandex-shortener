// Package main основная точка входа в программу
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/closable/go-yandex-shortener/internal/config"
	"github.com/closable/go-yandex-shortener/internal/handlers"
	"github.com/closable/go-yandex-shortener/internal/storage"
	"github.com/closable/go-yandex-shortener/internal/utils"
)

var buildVersion, buildDate, buildCommit = "N/A", "N/A", "N/A"

func main() {
	// go run -ldflags "-X main.buildVersion=v1.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" cmd/shortener/main.go
	// go build -ldflags "-X main.buildVersion=v1.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" cmd/shortener/main.go
	// start bin file -> ./main
	fmt.Printf("Build version:%s\nBuild date:%s\nBuild commit:%s\n", buildVersion, buildDate, buildCommit)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	idleConnsClosed := make(chan struct{})

	srv, isHTTPS := serverPrepare()
	go func() {
		<-sigint
		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		if err := srv.Shutdown(context.Background()); err != nil {
			// ошибки закрытия Listener
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	err := run(srv, isHTTPS) //   srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		// ошибки старта или остановки Listener
		panic(err)
	}

	<-idleConnsClosed
	fmt.Println("Server Shutdown gracefully")
}

// run старт сервиса
func run(srv *http.Server, isHTTPS bool) error {
	if isHTTPS {
		return srv.ListenAndServeTLS("", "")
	}
	return srv.ListenAndServe()
}

// serverPrepare подготовка и конфигурирование сервиса
func serverPrepare() (*http.Server, bool) {
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
		//os.Exit(1)
		panic(err)
	}

	handler := handlers.New(store, cfg.BaseURL, logger, 1)

	serverAddr, err := utils.MakeServerAddres(cfg.ServerAddress, cfg.EnableHTTPS)
	if err != nil {
		panic(err)
	}
	sugar.Infoln(storeMsg)
	sugar.Infoln("Running server on", serverAddr)

	if cfg.EnableHTTPS {
		server := &http.Server{
			Addr:    serverAddr,
			Handler: handler.InitRouter(),
			// для TLS-конфигурации используем менеджер сертификатов
			TLSConfig: utils.MekeTLS("closable@yandex.ru", "shortener").TLSConfig(),
		}
		return server, true

	} else {
		server := &http.Server{
			Addr:    serverAddr,
			Handler: handler.InitRouter(),
		}
		return server, false
	}

	//return http.ListenAndServe(cfg.ServerAddress, handler.InitRouter())
}

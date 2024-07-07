// Package main основная точка входа в программу
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/closable/go-yandex-shortener/internal/services"
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

	srv, err := services.New()
	if err != nil {
		panic(err)
	}

	go func() {
		<-sigint
		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		if err := srv.Shutdown(context.Background()); err != nil {
			// ошибки закрытия Listener
			msg := "HTTP(s) server Shutdown:"
			if srv.IsGRPC {
				msg = "gRPC server Shutdown:"
			}
			srv.Logger.Infoln(msg, err)
		}
		close(idleConnsClosed)
	}()

	err = srv.Run()
	if err != nil && err != srv.ErrClosed {
		// ошибки старта или остановки Listener
		panic(err)
	}

	<-idleConnsClosed
	srv.Logger.Infoln("Server Shutdown gracefully !")
}

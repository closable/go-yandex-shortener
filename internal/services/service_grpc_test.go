package services

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/closable/go-yandex-shortener/internal/handlers"
	pb "github.com/closable/go-yandex-shortener/internal/services/proto"
	"github.com/closable/go-yandex-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGRPChandlers(t *testing.T) {
	//var wg sync.WaitGroup
	logger := handlers.NewLogger()
	sugar := *logger.Sugar()
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}
	var store handlers.Storager
	store, _ = storage.NewMemory()
	s := grpc.NewServer()
	pb.RegisterShortenesServer(s, &ShortenerServer{
		store:   store,
		baseURL: "localhost:8080",
		logger:  &sugar,
	})

	go func() {
		// если делать правильно с обработкой ошики, то не знаю как решить проблему, корректного закрытия сервера
		// продолжает слушать даже при прямом указании остановиться
		s.Serve(listen)
	}()

	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewShortenesClient(conn)

	closer := func() {
		err := listen.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		fmt.Println("stop")
		s.Stop()
	}

	testingGetUrls(t, c)
	testingPing(t, c)
	testingStats(t, c)

	closer()

}

func testingStats(t *testing.T, c pb.ShortenesClient) {
	var req pb.GetStatsRequest
	res, err := c.GetStats(context.Background(), &req)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, res.Urls, "Wrong result")

}

func testingPing(t *testing.T, c pb.ShortenesClient) {
	var req pb.GetPingRequest
	res, err := c.Ping(context.Background(), &req)
	assert.NoError(t, err)
	assert.Equal(t, true, res.Ping, "Wrong result")

}

func testingGetUrls(t *testing.T, c pb.ShortenesClient) {
	var req pb.GetURLsRequest
	res, err := c.GetURLs(context.Background(), &req)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(res.Rows), "Wrong result")

	req1 := pb.GetShortenerRequest{
		UserId: 1, Url: "http://yandex.ru",
	}

	res1, err := c.GetShortener(context.Background(), &req1)
	assert.NoError(t, err)
	fmt.Println(res1)

	res, err = c.GetURLs(context.Background(), &req)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res.Rows), "Wrong result")
}

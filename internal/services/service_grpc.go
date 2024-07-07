// Package services use for configuration http(s)/gRPC services
package services

import (
	"context"

	"github.com/closable/go-yandex-shortener/internal/handlers"
	pb "github.com/closable/go-yandex-shortener/internal/services/proto"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ShortenerServer основная структура сервиса
type ShortenerServer struct {
	pb.UnimplementedShortenesServer
	store   handlers.Storager
	baseURL string
	logger  *zap.SugaredLogger
}

// GetShortener получение сокращения
func (s *ShortenerServer) GetShortener(ctx context.Context, in *pb.GetShortenerRequest) (*pb.GetShortenerResponse, error) {
	var response pb.GetShortenerResponse
	shorten, err := s.store.GetShortener(int(in.UserId), in.Url)
	response.Shortener = handlers.MakeShortenURL(shorten, s.baseURL)
	if err != nil {
		s.logger.Infoln("Shortener alredy present", shorten)
		return &response, status.Errorf(codes.Code(code.Code_ALREADY_EXISTS), "shortener is alredy exists")
	}
	s.logger.Infoln("Get shortener successfuly ", shorten)
	return &response, nil
}

// GetStats сбор статистики только в trasted подсети
func (s *ShortenerServer) GetStats(ctx context.Context, in *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	urls, users := s.store.GetStats()
	return &pb.GetStatsResponse{
		Urls: int32(urls), Users: int32(users),
	}, nil
}

// GetURLs список всех сокращений по указанному пользователю (только для memory store)
func (s *ShortenerServer) GetURLs(ctx context.Context, in *pb.GetURLsRequest) (*pb.GetURLsResponse, error) {
	var urls []*pb.UrlRow
	var response pb.GetURLsResponse
	u, err := s.store.GetURLs(int(in.UserId))
	if err != nil {
		return &response, err
	}

	for k, v := range u {
		urls = append(urls, &pb.UrlRow{
			OriginalURL: v, ShortURL: handlers.MakeShortenURL(k, s.baseURL),
		})
	}
	response = pb.GetURLsResponse{
		Rows: urls,
	}
	s.logger.Infoln("Get urls successfuly ", len(urls))
	return &response, nil
}

// Length количество записей в store
// func (s *ShortenerServer) Length(ctx context.Context, in *pb.CountURLsRequest) (*pb.CountURLsResponse, error) {
// 	response := pb.CountURLsResponse{
// 		Cnt: int32(s.store.Length()),
// 	}
// 	s.logger.Infoln("Count of elemets in store", s.store.Length())
// 	return &response, nil
// }

// Ping проверка работоспособности серевера
func (s *ShortenerServer) Ping(ctx context.Context, in *pb.GetPingRequest) (*pb.GetPingResponse, error) {
	return &pb.GetPingResponse{
		Ping: s.store.Ping(),
	}, nil
}

// SoftDeleteURLs заглушка
func (s *ShortenerServer) SoftDeleteURLs(userID int, key ...string) error {
	return nil
}

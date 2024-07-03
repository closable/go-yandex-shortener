package services

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/closable/go-yandex-shortener/internal/config"
	"github.com/closable/go-yandex-shortener/internal/handlers"
	pb "github.com/closable/go-yandex-shortener/internal/services/proto"
	"github.com/closable/go-yandex-shortener/internal/storage"
	"github.com/closable/go-yandex-shortener/internal/utils"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ServTransport структура фасада
type ServTransport struct {
	Store      handlers.Storager
	Logger     *zap.SugaredLogger
	HTTPServ   *http.Server
	GRPCServ   *grpc.Server
	GRPCListen net.Listener
	ServerAddr string
	IsHTTPS    bool
	IsGRPC     bool
	ErrClosed  error
}

// New новый экземпляр
func New() (*ServTransport, error) {
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
		sugar.Infoln("Error during configuring the protocol transport")
		return &ServTransport{}, err
	}
	sugar.Infoln(storeMsg)

	if cfg.UseGRPC {
		return GRPCconfigure(store, cfg.BaseURL, cfg.ServerAddress, cfg.TrastedSubnet, logger)
	} else {
		return HTTPconfigure(store, cfg.BaseURL, cfg.ServerAddress, cfg.TrastedSubnet, logger, cfg.EnableHTTPS)
	}
}

// Run service
func (t *ServTransport) Run() error {
	if t.IsGRPC {
		t.Logger.Infoln("Running gRPCS server on", t.ServerAddr)
		return t.GRPCServ.Serve(t.GRPCListen)
	} else {
		t.Logger.Infoln("Running server on", t.ServerAddr)
		if t.IsHTTPS {
			return t.HTTPServ.ListenAndServeTLS("", "")
		}
		return t.HTTPServ.ListenAndServe()
	}
}

// Shutdown use for Graceful shutdown
func (t *ServTransport) Shutdown(ctx context.Context) error {
	if t.IsGRPC {
		t.GRPCServ.GracefulStop()
		return nil
	}

	return t.HTTPServ.Shutdown(ctx)
}

// HTTPconfigure use for configure http(s) server
func HTTPconfigure(s handlers.Storager, baseURL, serverURL, trastedSubnet string, logger zap.Logger, https bool) (*ServTransport, error) {
	handler := handlers.New(s, baseURL, logger, 1, trastedSubnet)
	isHTTPS := false
	server := &http.Server{
		Addr:    serverURL,
		Handler: handler.InitRouter(),
	}

	if https {
		// для TLS-конфигурации используем менеджер сертификатов
		server.TLSConfig = utils.MekeTLS("closable@yandex.ru", "shortener").TLSConfig()
	}

	serverAddr, err := utils.MakeServerAddres(serverURL, https)
	if err != nil {
		panic(err)
	}
	sugar := *logger.Sugar()

	return &ServTransport{
		Store:      s,
		Logger:     &sugar,
		HTTPServ:   server,
		ServerAddr: serverAddr,
		IsHTTPS:    isHTTPS,
		IsGRPC:     false,
		ErrClosed:  http.ErrServerClosed,
	}, nil

}

// GRPCconfigure use for gRPC confugure server
func GRPCconfigure(s handlers.Storager, baseURL, serverURL, trastedSubnet string, logger zap.Logger) (*ServTransport, error) {

	//listen, err := net.Listen("tcp", ":3200")
	listen, err := net.Listen("tcp", serverURL)
	if err != nil {
		log.Fatal(err)
	}

	// создаём gRPC-сервер без зарегистрированной службы
	serv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(UnaryServerInterceptorOpts(trastedSubnet)),
	)
	sugar := *logger.Sugar()

	// регистрируем сервис
	pb.RegisterShortenesServer(serv, &ShortenerServer{
		store:   s,
		baseURL: serverURL,
		logger:  &sugar,
	})

	return &ServTransport{
		Store:      s,
		Logger:     &sugar,
		GRPCServ:   serv,
		GRPCListen: listen,
		ServerAddr: serverURL,
		IsHTTPS:    false,
		IsGRPC:     true,
		ErrClosed:  grpc.ErrServerStopped,
	}, nil

}

// UnaryServerInterceptorOpts middleware for trasted IP
func UnaryServerInterceptorOpts(trastedSubnet string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		method := strings.Split(info.FullMethod, "/")[2]

		if method == "GetStats" {
			trasted := false
			for k, v := range md {
				if strings.Contains(":authority x-real-ip x-forwarded-for", k) {

					// addrPort, err := netip.ParseAddrPort(v[0])
					ip, _, err := net.SplitHostPort(v[0])
					if err != nil {
						fmt.Println("!!!", err)
						continue
					}
					if strings.Contains(trastedSubnet, ip) {
						trasted = true
						break
					}
				}
			}

			if !trasted {
				return nil, status.Errorf(codes.Code(code.Code_CANCELLED), "Acccess only trasted iP")
			}
		}

		return handler(ctx, req)
	}
}

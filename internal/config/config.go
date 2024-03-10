package config

import (
	"flag"
	"strings"

	"github.com/caarlos0/env/v10"
)

type config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	FileStore     string `env:"FILE_STORAGE_PATH"`
	DSN           string `env:"DATABASE_DSN"`
}

var (
	FlagRunAddr   string
	FlagSendAddr  string
	FlagFileStore string
	FlagDSN       string
	configEnv     = config{}
)

func ParseConfigEnv() {
	env.Parse(&configEnv)
}

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	// адрес и порт куда отправлять сокращатель
	flag.StringVar(&FlagSendAddr, "b", "localhost:8080", "seneder address and port to run server")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.StringVar(&FlagFileStore, "f", "/tmp/short-url-db.json", "folder and path where to store data")
	flag.StringVar(&FlagDSN, "d", "postgres://postgres:1303@localhost:5432/postgres", "access to DBMS")

	flag.Parse()
}

// загружаем данные среды окружения
func LoadConfig() *config {
	ParseConfigEnv()
	ParseFlags()

	var config = &config{}

	if len(configEnv.ServerAddress) > 0 {
		config.ServerAddress = configEnv.ServerAddress
	}
	if len(configEnv.ServerAddress) == 0 && len(FlagRunAddr) > 0 {
		config.ServerAddress = FlagRunAddr
	}

	if len(configEnv.BaseURL) > 0 {
		config.BaseURL = configEnv.BaseURL
	}
	if len(configEnv.BaseURL) == 0 && len(FlagSendAddr) > 0 {
		config.BaseURL = FlagSendAddr
	}

	if len(configEnv.FileStore) > 0 {
		config.FileStore = configEnv.FileStore
	}
	if len(configEnv.FileStore) == 0 && len(FlagFileStore) > 0 {
		config.FileStore = FlagFileStore
	}

	if len(configEnv.DSN) > 0 {
		config.DSN = configEnv.DSN
	}
	if len(configEnv.DSN) == 0 && len(FlagDSN) > 0 {
		config.DSN = FlagDSN
	}
	if !strings.Contains(config.DSN, "sslmode") {
		config.DSN += "?sslmode=disable"
	}

	return config
}

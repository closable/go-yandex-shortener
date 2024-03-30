package config

import (
	"flag"
	"fmt"

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
	//flag.StringVar(&FlagFileStore, "f", "/tmp/short-url-db.json", "folder and path where to store data")
	flag.StringVar(&FlagFileStore, "f", "", "folder and path where to store data")
	// flag.StringVar(&FlagDSN, "d", "postgres://postgres:1303@localhost:5432/postgres", "access to DBMS")
	flag.StringVar(&FlagDSN, "d", "", "access to DBMS")

	flag.Parse()
}

// загружаем данные среды окружения
func LoadConfig() *config {
	ParseConfigEnv()
	ParseFlags()

	var config = &config{}
	config.ServerAddress = firstValue(&configEnv.ServerAddress, &FlagRunAddr)
	config.BaseURL = firstValue(&configEnv.BaseURL, &FlagSendAddr)
	config.FileStore = firstValue(&configEnv.FileStore, &FlagFileStore)
	config.DSN = firstValue(&configEnv.DSN, &FlagDSN)
	// if !strings.Contains(config.DSN, "sslmode") {
	// 	config.DSN = fmt.Sprintf("%s?sslmode=disable", config.DSN)
	// }
	fmt.Println(configEnv.DSN, FlagDSN)
	return config
}

func firstValue(valEnv *string, valFlag *string) string {
	if len(*valEnv) > 0 {
		return *valEnv
	}
	return *valFlag
}

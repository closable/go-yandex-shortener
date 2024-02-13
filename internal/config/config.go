package config

import (
	"flag"

	"github.com/caarlos0/env/v10"
)

type config struct {
	//ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	ServerAddress string `env:"SERVER_ADDRESS"`
	//BaseURL       string `env:"BASE_URL" envDefault:"localhost:8080"`
	BaseURL string `env:"BASE_URL"`
}

var ConfigEnv = config{}

func ParseConfigEnv() {
	env.Parse(&ConfigEnv)
}

func GetEnvParam(typeVar string) string {

	if typeVar == "RUN_SERVER" {
		if len(ConfigEnv.ServerAddress) > 0 {
			return ConfigEnv.ServerAddress
		}
		if len(FlagRunAddr) > 0 {
			return FlagRunAddr
		}
	}

	if typeVar == "SND_SERVER" {
		if len(ConfigEnv.BaseURL) > 0 {
			return ConfigEnv.BaseURL
		}
		if len(FlagSendAddr) > 0 {
			return FlagSendAddr
		}
	}
	return "localhost:8080"
}

// экспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var FlagRunAddr string
var FlagSendAddr string

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	// адрес и порт куда отправлять сокращатель
	flag.StringVar(&FlagSendAddr, "b", "localhost:8080", "seneder address and port to run server")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}

// загружаем данные среды окружения
func LoadConfig() {
	ParseConfigEnv()
	ParseFlags()
}

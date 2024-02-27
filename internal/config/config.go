package config

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v10"
)

type config struct {
	//ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	ServerAddress string `env:"SERVER_ADDRESS"`
	//BaseURL       string `env:"BASE_URL" envDefault:"localhost:8080"`
	BaseURL   string `env:"BASE_URL"`
	FileStore string `env:"FILE_STORAGE_PATH"`
}

var (
	FlagRunAddr   string
	FlagSendAddr  string
	FlagFileStore string
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

	//flag.StringVar(&FlagFileStore, "f", "./tmp/YyHvN0A", "folder and path where to store data")

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
	if len(FlagRunAddr) > 0 {
		config.ServerAddress = FlagRunAddr
	}

	if len(configEnv.BaseURL) > 0 {
		config.BaseURL = configEnv.BaseURL
	}
	if len(FlagSendAddr) > 0 {
		config.BaseURL = FlagSendAddr
	}

	if len(configEnv.FileStore) > 0 {
		config.FileStore = configEnv.FileStore
	}
	if len(FlagFileStore) > 0 {
		config.FileStore = FlagFileStore
	}

	// if len(config.FileStore) > 0 {
	// 	// for UNIX the /tmp folder is usually there, but it needs to be corrected relative to the working directory
	// 	fileNameCorrected := fmt.Sprintf(".%s", FlagFileStore)
	// 	CreateNotIxistingFolders(fileNameCorrected)
	// 	config.FileStore = fileNameCorrected
	// }

	return config
}

func CreateNotIxistingFolders(fileName string) {
	if _, err := os.Stat(fileName); err != nil {
		path := filepath.Dir(fileName)
		os.MkdirAll(path, os.ModePerm)
	}
}

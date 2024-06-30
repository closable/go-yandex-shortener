// Package config служит для получения данных от входящего окружения
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/caarlos0/env/v10"
)

type fileConfig struct {
	// аналог переменной окружения SERVER_ADDRESS или флага -a
	ServerAddress string `json:"server_address"`
	// аналог переменной окружения BASE_URL или флага -b
	BaseURL string `json:"base_url"`
	// аналог переменной окружения FILE_STORAGE_PATH или флага -f
	FileStore string `json:"file_storage_path"`
	// аналог переменной окружения DATABASE_DSN или флага -d
	DSN string `json:"database_dsn"`
	// аналог переменной окружения ENABLE_HTTPS или флага -s
	EnableHTTPS bool `json:"enable_https"`
	//  аналог переменной окружения TRUSTED_SUBNET
	TrastedSubnet string `json:"trusted_subnet"`
}

// config описание структур данных среды окружения
type config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	FileStore     string `env:"FILE_STORAGE_PATH"`
	DSN           string `env:"DATABASE_DSN"`
	EnableHTTPS   bool   `env:"ENABLE_HTTPS"`
	UseConfig     bool   `env:"CONFIG"`
	TrastedSubnet string `env:"TRUSTED_SUBNET"`
}

// переменные
var (
	// Адрес сервера
	FlagRunAddr string
	// Адрес выдачи информации
	FlagSendAddr string
	// Использование файлового хранилища
	FlagFileStore string
	// Использование СУБД
	FlagDSN string
	// Активация https
	FlagHTTPS bool
	// Использовать config.json
	FlagConfig bool
	// Доверенная подсеть
	FlagTrasted string
	configEnv   = config{}
)

// ParseConfigEnv парсинг переменных среды окружения
func ParseConfigEnv() {
	env.Parse(&configEnv)
}

// ParseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	// адрес и порт куда отправлять сокращатель
	flag.StringVar(&FlagSendAddr, "b", "localhost:8080", "sender address and port to run server")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	//flag.StringVar(&FlagFileStore, "f", "/tmp/short-url-db.json", "folder and path where to store data")
	flag.StringVar(&FlagFileStore, "f", "", "folder and path where to store data")
	//flag.StringVar(&FlagDSN, "d", "postgres://postgres:1303@localhost:5432/postgres", "access to DBMS")
	flag.StringVar(&FlagDSN, "d", "", "access to DBMS")
	flag.BoolVar(&FlagHTTPS, "s", false, "use https")
	flag.BoolVar(&FlagConfig, "c", false, "use config file")
	flag.StringVar(&FlagTrasted, "t", "", "use trasted subnet")

	flag.Parse()
}

// LoadConfig загружаем данные среды окружения
func LoadConfig() *config {
	ParseConfigEnv()
	ParseFlags()

	var config = &config{}

	config.ServerAddress = firstValue(&configEnv.ServerAddress, &FlagRunAddr)
	config.BaseURL = firstValue(&configEnv.BaseURL, &FlagSendAddr)
	config.FileStore = firstValue(&configEnv.FileStore, &FlagFileStore)
	config.DSN = firstValue(&configEnv.DSN, &FlagDSN)
	config.EnableHTTPS = configEnv.EnableHTTPS || FlagHTTPS
	config.UseConfig = configEnv.UseConfig || FlagConfig
	config.TrastedSubnet = firstValue(&configEnv.TrastedSubnet, &FlagTrasted)

	if config.UseConfig {
		updateFromConfig(config)
	}

	return config
}

// firstValue вспомогательная функция для выбора входящих значений
func firstValue(valEnv *string, valFlag *string) string {
	if len(*valEnv) > 0 {
		return *valEnv
	}
	return *valFlag
}

// updateFromCnfig обновление переменных из файла сонфигурации
func updateFromConfig(c *config) error {
	// f, err := os.OpenFile("../../config.json", os.O_RDONLY|os.O_RDWR, 0755)
	f, err := os.OpenFile("config.json", os.O_RDONLY|os.O_RDWR, 0755)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	res := &fileConfig{}
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil
	}

	c.ServerAddress = firstValue(&res.ServerAddress, &c.ServerAddress)
	c.BaseURL = firstValue(&res.BaseURL, &c.BaseURL)
	c.FileStore = firstValue(&res.FileStore, &c.FileStore)
	c.DSN = firstValue(&res.DSN, &c.DSN)
	c.EnableHTTPS = res.EnableHTTPS || c.EnableHTTPS
	c.TrastedSubnet = firstValue(&res.TrastedSubnet, &c.TrastedSubnet)
	return nil
}

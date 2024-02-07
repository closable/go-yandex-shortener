package config

import (
	"flag"
)

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

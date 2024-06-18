// Package utils реализует вспомогательные функции для аботы приложения
package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

// Описание структур данных
type (
	// BatchBody структура описания тела для массового импорта данных
	BatchBody struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}
)

// ValidateURL check url
func ValidateURL(txtURL string) bool {
	u, err := url.Parse(txtURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// GetRandomKey вспомогательная функция для получения провольного ключа
func GetRandomKey(length int) string {
	chars := "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	key := ""
	for i := 0; i < length; i++ {
		c := chars[rand.Intn(len(chars))]
		key += string(c)

	}
	return key
}

// GenerateBatchBody вспомогательная функция для генерации тела массовой загрузки
func GenerateBatchBody(itemsCnt int) {
	arr := []BatchBody{}

	for i := 0; i <= itemsCnt; i++ {

		arr = append(arr, BatchBody{
			CorrelationID: GetRandomKey(5),
			OriginalURL:   fmt.Sprintf("http://%s/yandex.ru/%s", GetRandomKey(4), GetRandomKey(7)),
		})
	}

	resp, err := json.MarshalIndent(arr, " ", "    ")
	if err != nil {
		fmt.Sprintln(err)
	}

	file, err := os.OpenFile("/tmp/batch_body.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Sprintln(err)
	}
	defer file.Close()

	file.Write(resp)

	//fmt.Println(string(resp))
}

// MakeServerAddres составляет корректный адрес сервера в зависимости от флага https
func MakeServerAddres(addr string, flagHTTPS string) (string, error) {
	if len(flagHTTPS) == 0 {
		return addr, nil
	}

	u, err := url.Parse(addr)
	if err != nil {
		return "", err
	}
	if u.Host == "localhost" || u.Host == "127.0.0.1" {
		return ":443", nil
	} else {
		return fmt.Sprintf("%s:443", u.Host), nil
	}
}

func MekeTLS(email string, whiteList string) *autocert.Manager {
	manager := &autocert.Manager{
		// директория для хранения сертификатов
		Cache: autocert.DirCache("cache-dir"),
		// функция, принимающая Terms of Service издателя сертификатов
		Prompt: autocert.AcceptTOS,
		Email:  email,
		// перечень доменов, для которых будут поддерживаться сертификаты
		HostPolicy: autocert.HostWhitelist(whiteList, fmt.Sprintf("www.%s", whiteList)),
	}
	// конструируем сервер с поддержкой TLS
	return manager
}

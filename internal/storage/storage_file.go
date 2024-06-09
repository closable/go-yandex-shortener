package storage

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/closable/go-yandex-shortener/internal/utils"
)

type (
	//Event Структора для работы JSON
	Event struct {
		UUID       string `json:"uud"`
		ShortURL   string `json:"short_url"`
		OriginlURL string `json:"original_url"`
	}
	// StoregeFile структура файлового хранилища
	StoregeFile struct {
		File    *os.File
		Encoder *json.Encoder
	}
	// Consumer Структура для работы с файлововым хранилищем
	Consumer struct {
		File    *os.File
		Decoder *json.Decoder
	}
)

// NewFile новый экземпляр хранения
func NewFile(fileName string) (*StoregeFile, error) {
	os.MkdirAll(filepath.Dir(fileName), os.ModePerm)

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &StoregeFile{
		File:    file,
		Encoder: json.NewEncoder(file),
	}, nil

}

// WriteEvent для записи в файловое хранилище
func (s *StoregeFile) WriteEvent(event *Event) error {
	return s.Encoder.Encode(&event)
}

// NewConsumer создание нового экземпляра
func NewConsumer(fileName string) (*Consumer, error) {

	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		File:    file,
		Decoder: json.NewDecoder(file),
	}, nil
}

// Close функция хакрытия файового хранилища
func (s *StoregeFile) Close() error {
	return s.File.Close()
}

// func loadDataFromFile(file *os.File) error {
// 	body, err := io.ReadAll(file)
// 	if err != nil {
// 		return err
// 	}

// 	rows := strings.Split(string(body), "\n")

// 	for _, v := range rows {
// 		event := &Event{}
// 		err := json.Unmarshal([]byte(v), event)
// 		if err != nil {
// 			file.Close()
// 			break
// 		}
// 		//st.AddItem(event.ShortURL, event.OriginlURL)
// 	}
// 	return nil
// }

// FindKeyByValue поиск ключа по значению
func (s *StoregeFile) FindKeyByValue(urlText string) string {

	consumer, err := NewConsumer(s.File.Name())
	if err != nil {
		return ""
	}

	defer consumer.File.Close()

	body, _ := io.ReadAll(consumer.File)

	rows := strings.Split(string(body), "\n")

	for _, v := range rows {
		event := &Event{}

		err := json.Unmarshal([]byte(v), event)
		if err != nil {
			consumer.File.Close()
			break
		}

		if urlText == event.OriginlURL {
			return event.ShortURL
		}
	}
	return ""
}

// FindExistingKey поиск существующего ключа
func (s *StoregeFile) FindExistingKey(keyText string) (string, bool) {

	consumer, err := NewConsumer(s.File.Name())
	if err != nil {
		return "", false
	}
	defer consumer.File.Close()

	body, err := io.ReadAll(consumer.File)
	if err != nil {
		return "", false
	}

	rows := strings.Split(string(body), "\n")

	for _, v := range rows {
		event := &Event{}
		err := json.Unmarshal([]byte(v), event)
		if err != nil {
			consumer.File.Close()
			break
		}

		if keyText == event.ShortURL {
			return event.OriginlURL, true
		}
	}
	return "", false
}

// GetShortener получение сокращения
func (s *StoregeFile) GetShortener(userID int, urlText string) (string, error) {
	var shorterner string

	shorterner = s.FindKeyByValue(urlText)
	if len(shorterner) > 0 {
		return shorterner, nil
	}

	shorterner = utils.GetRandomKey(6)
	err := s.WriteEvent(&Event{
		UUID:       utils.GetRandomKey(10),
		ShortURL:   shorterner,
		OriginlURL: string(urlText),
	})

	if err != nil {
		return "", err
	}
	return shorterner, nil
}

// Ping healthcheck routine
func (s *StoregeFile) Ping() bool {
	shortener, _ := s.GetShortener(0, "ping")
	return len(shortener) != 0
}

// PrepareStore заглушка для удовлетворения интерфейсу
func (s *StoregeFile) PrepareStore() {
}

// GetURLs заглушка для удовлетворения интерфейсу
func (s *StoregeFile) GetURLs(userID int) (map[string]string, error) {
	var result = make(map[string]string)
	return result, nil
}

// SoftDeleteURLs заглушка для удовлетворения интерфейсу
func (s *StoregeFile) SoftDeleteURLs(userID int, key ...string) error {
	return nil
}

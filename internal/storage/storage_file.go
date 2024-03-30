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
	Event struct {
		UUID       string `json:"uud"`
		ShortURL   string `json:"short_url"`
		OriginlURL string `json:"original_url"`
	}
	StoregeFile struct {
		File    *os.File
		Encoder *json.Encoder
	}
	Consumer struct {
		File    *os.File
		Decoder *json.Decoder
	}
)

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

func (st *StoregeFile) WriteEvent(event *Event) error {
	return st.Encoder.Encode(&event)
}

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

func (st *StoregeFile) Close() error {
	return st.File.Close()
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

func (st *StoregeFile) FindKeyByValue(urlText string) string {

	consumer, err := NewConsumer(st.File.Name())
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

func (st *StoregeFile) FindExistingKey(keyText string) (string, bool) {

	consumer, err := NewConsumer(st.File.Name())
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

func (st *StoregeFile) GetShortener(urlText string) (string, error) {
	var shorterner string

	shorterner = st.FindKeyByValue(urlText)
	if len(shorterner) > 0 {
		return shorterner, nil
	}

	shorterner = utils.GetRandomKey(6)
	err := st.WriteEvent(&Event{
		UUID:       utils.GetRandomKey(10),
		ShortURL:   shorterner,
		OriginlURL: string(urlText),
	})

	if err != nil {
		return "", err
	}
	return shorterner, nil
}

func (st *StoregeFile) Ping() bool {
	shortener, _ := st.GetShortener("ping")
	return len(shortener) != 0
}

func (st *StoregeFile) PrepareStore() {
}

func (st *StoregeFile) GetURLs(userID int) (map[string]string, error) {
	var result = make(map[string]string)
	return result, nil
}

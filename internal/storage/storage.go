// Package storage реализует функцонал хранения данных
package storage

import (
	"errors"
	"strings"
	"sync"

	"github.com/closable/go-yandex-shortener/internal/utils"
)

// Store структураописания модели ханения в памяти
type Store struct {
	mu   sync.Mutex
	Urls map[string]string
}

// NewMemory Создание нового экземпляра хранения в памяти
func NewMemory() (*Store, error) {
	return &Store{Urls: make(map[string]string)}, nil
}

// Length вспомогательная функция
func (s *Store) Length() int {
	return len(s.Urls)
}

// AddItem Добавление нового элемента
func (s *Store) AddItem(key string, url string) (string, error) {
	s.Urls[key] = url
	return key, nil
}

// GetShortener полчение сокращения URL
func (s *Store) GetShortener(userID int, txtURL string) (string, error) {
	shortener := ""
	// it needs for exclude existing urls
	s.mu.Lock()
	defer s.mu.Unlock()
	key := s.FindKeyByValue(string(txtURL))

	if len(key) == 0 {
		shortener = utils.GetRandomKey(6)
		for {
			// exclude existing keys
			_, ok := s.FindExistingKey(shortener)
			if !ok {
				break
			}
		}
		s.Urls[shortener] = txtURL
	} else {
		return key, errors.New("409")
	}

	return shortener, nil
}

// FindKeyByValue Поиск ключа по значениб
func (s *Store) FindKeyByValue(urlText string) string {
	for key, value := range s.Urls {
		if strings.Contains(value, urlText) {
			return key
		}
	}
	return ""
}

// FindExistingKey поиск существующего ключа
func (s *Store) FindExistingKey(keyText string) (string, bool) {

	value, ok := s.Urls[keyText]

	return value, ok
}

// Ping вспомогательня функция
func (s *Store) Ping() bool {
	s.Urls["ping"] = "ping"
	return s.Urls["ping"] == "ping"
}

// PrepareStore Заглушка для удовлетворения интерфейсу
func (s *Store) PrepareStore() {
}

// GetURLs Заглушка для удовлетворения интерфейсу
func (s *Store) GetURLs(userID int) (map[string]string, error) {
	return s.Urls, nil
}

// SoftDeleteURLs Заглушка для удовлетворения интерфейсу
func (s *Store) SoftDeleteURLs(userID int, key ...string) error {
	return nil
}

// GetStats сбор данных для совместимости с ограниченным функционалом
func (s *Store) GetStats() (int, int) {
	return s.Length(), 0
}

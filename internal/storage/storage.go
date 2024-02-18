package storage

import (
	"sync"

	"github.com/closable/go-yandex-shortener/internal/utils"
)

// var Storage = struct {
// 	mu   sync.Mutex
// 	Urls map[string]string
// }{Urls: make(map[string]string)}

type Store struct {
	mu   sync.Mutex
	Urls map[string]string
}

func New() *Store {
	return &Store{Urls: make(map[string]string)}
}

func (s *Store) GetShortener(txtURL string) string {
	shortener := ""
	// it needs for exclude existing urls
	s.mu.Lock()
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
		shortener = key
	}

	s.mu.Unlock()
	return shortener
}

func (s *Store) FindKeyByValue(urlText string) string {
	for key, value := range s.Urls {
		if value == urlText {
			return key
		}
	}
	return ""
}

func (s *Store) FindExistingKey(keyText string) (string, bool) {

	value, ok := s.Urls[keyText]

	return value, ok
}

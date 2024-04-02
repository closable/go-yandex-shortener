package storage

import (
	"errors"
	"strings"
	"sync"

	"github.com/closable/go-yandex-shortener/internal/utils"
)

type Store struct {
	mu   sync.Mutex
	Urls map[string]string
}

func NewMemory() (*Store, error) {
	return &Store{Urls: make(map[string]string)}, nil
}

func (s *Store) Length() int {
	return len(s.Urls)
}

func (s *Store) AddItem(key string, url string) (string, error) {
	s.Urls[key] = url
	return key, nil
}

func (s *Store) GetShortener(txtURL string) (string, error) {
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

func (s *Store) FindKeyByValue(urlText string) string {
	for key, value := range s.Urls {
		if strings.Contains(value, urlText) {
			return key
		}
	}
	return ""
}

func (s *Store) FindExistingKey(keyText string) (string, bool) {

	value, ok := s.Urls[keyText]

	return value, ok
}

func (s *Store) Ping() bool {
	s.Urls["ping"] = "ping"
	return s.Urls["ping"] == "ping"
}

func (s *Store) PrepareStore() {
}

func (s *Store) GetURLs(userID int) (map[string]string, error) {
	return s.Urls, nil
}

func (s *Store) SoftDeleteURLs(userID int, key ...string) error {
	return nil
}

package storage

import (
	"math/rand"
	"sync"
)

var Storage = struct {
	mu   sync.Mutex
	Urls map[string]string
}{Urls: make(map[string]string)}

func GetShortener(txtURL string) string {
	shortener := ""
	// it needs for exclude existing urls
	Storage.mu.Lock()
	key := FindKeyByValue(string(txtURL))

	if len(key) == 0 {
		shortener = GetShortnerKey(6)
		Storage.Urls[shortener] = txtURL
	} else {
		shortener = key
	}

	Storage.mu.Unlock()
	return shortener
}

func GetShortnerKey(length int) string {
	chars := "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	shortener := ""
	for {
		for i := 0; i < length; i++ {
			c := chars[rand.Intn(len(chars))]
			shortener += string(c)

		}
		// exclude existing keys
		if !(FindExistingKey(shortener)) {
			return shortener
		}
	}
}

func FindKeyByValue(urlText string) string {
	for key, value := range Storage.Urls {
		if value == urlText {
			return key
		}
	}
	return ""
}

func FindExistingKey(keyText string) bool {

	_, ok := Storage.Urls[keyText]

	return ok
}

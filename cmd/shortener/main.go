package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

// var Storage = make(map[string]string)

var Storage = struct {
	mu   sync.Mutex
	Urls map[string]string
}{Urls: make(map[string]string)}

func main() {

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(`:8080`, http.HandlerFunc(webhook))
}

func webhook(w http.ResponseWriter, r *http.Request) {

	body := ""
	shortener := ""

	w.Header().Set("Content-Type", "plain/text")
	switch method := r.Method; method {
	case http.MethodPost:

		info, err := io.ReadAll(r.Body)
		if err != nil || len(info) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error! Request body is empty!"))
			return
		}

		if !(validateURL(string(info))) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error! Check url it's should seems as like this 'http[s]://example.com'"))
			return
		}

		shortener = getShortener(string(info))

		w.WriteHeader(http.StatusCreated)
		body = fmt.Sprintf("http://%s/%s", r.Host, shortener)

	case http.MethodGet:
		path := r.URL.Path
		if len(path) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error! id is empty!"))
			return

		}
		shortener = path[1:]

		if !findExistingKey(shortener) || len(shortener) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error! id is not found or empty"))
			return
		}

		w.Header().Set("Location", Storage.Urls[shortener])
		w.WriteHeader(http.StatusTemporaryRedirect)
		// body = fmt.Sprintf("%s ", Storage.Urls[shortener])

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Error! Please, check documentation again!"))
		return
	}

	w.Write([]byte(body))

}

func getShortener(txtURL string) string {
	shortener := ""
	// it needs for exclude existing urls
	Storage.mu.Lock()
	key := findKeyByValue(string(txtURL))

	if len(key) == 0 {
		shortener = getShortnerKey(6)
		Storage.Urls[shortener] = txtURL
	} else {
		shortener = key
	}

	Storage.mu.Unlock()
	return shortener
}

// check url
func validateURL(txtURL string) bool {
	u, err := url.Parse(txtURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func getShortnerKey(length int) string {
	chars := "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	shortener := ""
	for {
		for i := 0; i < length; i++ {
			c := chars[rand.Intn(len(chars))]
			shortener += string(c)

		}
		// exclude existing keys
		if !(findExistingKey(shortener)) {
			return shortener
		}
	}
}

func findKeyByValue(urlText string) string {
	for key, value := range Storage.Urls {
		if value == urlText {
			return key
		}
	}
	return ""
}

func findExistingKey(keyText string) bool {

	_, ok := Storage.Urls[keyText]

	return ok

	/*for key := range Storage.Urls {
		if key == keyText {
			return true
		}
	}
	return false*/
}

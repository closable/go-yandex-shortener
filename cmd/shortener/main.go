package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
)

var Storage = make(map[string]string)

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

		if !(validateUrl(string(info))) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error! Check url it's should seems as like this 'http[s]://example.com'"))
			return
		}

		shortener = getShortener(string(info))

		w.WriteHeader(http.StatusCreated)
		body = fmt.Sprintf("%s/%s", r.Host, shortener)
		//fmt.Println(len(Storage))
		//body += fmt.Sprintf("url: %s \r\n short: %s ", string(info), shortener)

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

		w.WriteHeader(http.StatusTemporaryRedirect)
		body = fmt.Sprintf("%s ", Storage[shortener])

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Error! Please, check documentation again!"))
		return
	}

	// fmt.Println(Storage)

	w.Write([]byte(body))

}

func getShortener(urlText string) string {
	shortener := ""
	// it needs for exclude existing urls
	key := findKeyByValue(string(urlText))

	if len(key) == 0 {
		shortener = getShortnerKey(6)
		Storage[shortener] = urlText
	} else {
		shortener = key
	}
	return shortener
}

// check url
func validateUrl(txtUrl string) bool {
	u, err := url.Parse(txtUrl)
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
	for key, value := range Storage {
		if value == urlText {
			return key
		}
	}
	return ""
}

func findExistingKey(keyText string) bool {
	for key := range Storage {
		if key == keyText {
			return true
		}
	}
	return false
}

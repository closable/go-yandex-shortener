package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/closable/go-yandex-shortener/internal/storage"
	"github.com/closable/go-yandex-shortener/internal/utils"
)

func Webhook(w http.ResponseWriter, r *http.Request) {

	switch method := r.Method; method {
	case http.MethodPost:
		GenerateShortener(w, r)

	case http.MethodGet:
		GetEndpointByShortener(w, r)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Error! Please, check documentation again!"))
		return
	}
}

func GenerateShortener(w http.ResponseWriter, r *http.Request) {

	shortener := ""

	info, err := io.ReadAll(r.Body)
	if err != nil || len(info) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error! Request body is empty!"))
		return
	}

	if !(utils.ValidateURL(string(info))) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error! Check url it's should seems as like this 'http[s]://example.com'"))
		return
	}

	shortener = storage.GetShortener(string(info))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	body := fmt.Sprintf("http://%s/%s", r.Host, shortener)
	w.Write([]byte(body))

}

func GetEndpointByShortener(w http.ResponseWriter, r *http.Request) {
	shortener := ""
	path := r.URL.Path
	if len(path) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error! id is empty!"))
		return

	}
	shortener = path[1:]

	if !storage.FindExistingKey(shortener) || len(shortener) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error! id is not found or empty"))
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", storage.Storage.Urls[shortener])
	w.WriteHeader(http.StatusTemporaryRedirect)
}

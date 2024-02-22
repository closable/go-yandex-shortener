package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/closable/go-yandex-shortener/internal/utils"
)

type Storager interface {
	GetShortener(txtURL string) string
	FindExistingKey(keyText string) (string, bool)
}

type URLHandler struct {
	store   Storager
	baseURL string
}

func New(st Storager, baseURL string) *URLHandler {
	return &URLHandler{
		store:   st,
		baseURL: baseURL,
	}
}

func (h *URLHandler) GenerateShortener(w http.ResponseWriter, r *http.Request) {

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
	shortener = h.store.GetShortener(string(info))

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	adr, _ := url.Parse(h.baseURL)

	body := ""
	if len(adr.Host) == 0 {
		body = fmt.Sprintf("http://%s/%s", h.baseURL, shortener)
	} else {
		body = fmt.Sprintf("%s/%s", h.baseURL, shortener)
	}

	w.Write([]byte(body))

}

func (h *URLHandler) GetEndpointByShortener(w http.ResponseWriter, r *http.Request) {
	shortener := ""
	path := r.URL.Path

	if len(path) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error! id is empty!"))
		return

	}
	shortener = path[1:]
	url, ok := h.store.FindExistingKey(shortener)

	if !ok || len(shortener) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error! id is not found or empty"))
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

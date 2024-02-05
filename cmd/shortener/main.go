package main

import (
	"net/http"

	"github.com/closable/go-yandex-shortener/internal/handlers"
)

func main() {

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(`:8080`, http.HandlerFunc(handlers.Webhook))
}

/*func webhook(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "plain/text")
	switch method := r.Method; method {
	case http.MethodPost:
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

		w.WriteHeader(http.StatusCreated)
		body := fmt.Sprintf("http://%s/%s", r.Host, shortener)
		w.Write([]byte(body))

	case http.MethodGet:
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

		w.Header().Set("Location", storage.Storage.Urls[shortener])
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Error! Please, check documentation again!"))
		return
	}
}*/

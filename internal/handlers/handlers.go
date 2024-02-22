package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/closable/go-yandex-shortener/internal/utils"
	"go.uber.org/zap"
)

type Storager interface {
	GetShortener(txtURL string) string
	FindExistingKey(keyText string) (string, bool)
}

type URLHandler struct {
	store   Storager
	baseURL string
	logger  zap.Logger
}

var (
	errBody    = "Error! Request body is empty!"
	errURL     = "Error! Check url it's should seems as like this 'http[s]://example.com'"
	emptyId    = "Error! id is empty!"
	notFoundId = "Error! id is not found or empty"
)

func New(st Storager, baseURL string, logger zap.Logger) *URLHandler {
	return &URLHandler{
		store:   st,
		baseURL: baseURL,
		logger:  logger,
	}
}

type ResponseWriter interface {
	//Header() Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func (h *URLHandler) GenerateShortener(w http.ResponseWriter, r *http.Request) {
	sugar := *h.logger.Sugar()
	shortener := ""

	info, err := io.ReadAll(r.Body)
	if err != nil || len(info) == 0 {

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errBody))
		sugar.Debugln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", errBody,
		)
		return
	}

	if !(utils.ValidateURL(string(info))) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errURL))
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", errURL,
		)
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
	sugar := *h.logger.Sugar()
	shortener := ""
	path := r.URL.Path

	if len(path) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(emptyId))
		sugar.Debugln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", emptyId,
		)
		return

	}
	shortener = path[1:]
	url, ok := h.store.FindExistingKey(shortener)

	if !ok || len(shortener) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(notFoundId))
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", notFoundId,
		)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

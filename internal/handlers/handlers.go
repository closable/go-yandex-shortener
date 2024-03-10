package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/closable/go-yandex-shortener/internal/storage"
	"github.com/closable/go-yandex-shortener/internal/utils"
	"go.uber.org/zap"
)

type Storager interface {
	GetShortener(txtURL string) string
	FindExistingKey(keyText string) (string, bool)
	Length() int
	AddItem(key string, url string)
}

type (
	URLHandler struct {
		store     Storager
		baseURL   string
		logger    zap.Logger
		fileStore string
		dbms      *storage.StoreDBMS
		maxLength int64
	}
	JSONRequest struct {
		URL string `json:"url"`
	}
	JSONRespond struct {
		Result string `json:"result"`
	}
)

var (
	errBody    = "Error! Request body is empty!"
	errURL     = "Error! Check url it's should seems as like this 'http[s]://example.com'"
	emptyID    = "Error! id is empty!"
	notFoundID = "Error! id is not found or empty"
)

func New(st Storager, baseURL string, logger zap.Logger, fileStore string, dbms *storage.StoreDBMS, maxLength int64) *URLHandler {
	// load stored data
	if len(fileStore) > 0 {
		consumer, err := NewConsumer(fileStore)
		if err != nil {
			logger.Fatal("File not found")
		}
		defer consumer.file.Close()
		loadDataFromFile(st, consumer.file)
	}

	return &URLHandler{
		store:     st,
		baseURL:   baseURL,
		logger:    logger,
		fileStore: fileStore,
		dbms:      dbms,
		maxLength: maxLength, // will compress if content-length > maxLength
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

func loadDataFromFile(st Storager, file *os.File) error {
	body, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	rows := strings.Split(string(body), "\n")

	for _, v := range rows {
		event := &Event{}
		err := json.Unmarshal([]byte(v), event)
		if err != nil {
			file.Close()
			break
		}
		st.AddItem(event.ShortURL, event.OriginlURL)
	}
	return nil
}

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

func (uh *URLHandler) GenerateShortener(w http.ResponseWriter, r *http.Request) {
	sugar := *uh.logger.Sugar()
	shortener := ""
	var reader io.Reader

	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}

	info, err := io.ReadAll(reader)
	if err != nil || len(info) == 0 {

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errBody))
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", errBody,
		)
		return
	}

	if !(utils.ValidateURL(string(info))) {
		// change behaviour when requeust doesn't have the protocol 26-02-24
		info = []byte(fmt.Sprintf("http://%s", info))
	}
	shortener = uh.store.GetShortener(string(info))

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	adr, _ := url.Parse(uh.baseURL)

	body := ""
	if len(adr.Host) == 0 {
		body = fmt.Sprintf("http://%s/%s", uh.baseURL, shortener)
	} else {
		body = fmt.Sprintf("%s/%s", uh.baseURL, shortener)
	}

	if len(uh.fileStore) > 0 {
		producer, err := NewProducer(uh.fileStore)
		if err != nil {
			sugar.Fatal(err)
		}

		defer producer.Close()
		if err := producer.WriteEvent(&Event{
			UUID:       uint(uh.store.Length()),
			ShortURL:   shortener,
			OriginlURL: string(info),
		}); err != nil {
			sugar.Fatal(err)
		}
	}

	w.Write([]byte(body))

}

func (uh *URLHandler) GetEndpointByShortener(w http.ResponseWriter, r *http.Request) {
	sugar := *uh.logger.Sugar()
	shortener := ""
	path := r.URL.Path

	if len(path) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(emptyID))
		sugar.Debugln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", emptyID,
		)
		return

	}
	shortener = path[1:]
	url, ok := uh.store.FindExistingKey(shortener)

	if !ok || len(shortener) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(notFoundID))
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", notFoundID,
		)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (uh *URLHandler) GenerateJSONShortener(w http.ResponseWriter, r *http.Request) {
	sugar := *uh.logger.Sugar()
	shortener := ""
	body := ""
	var jsonURL = &JSONRequest{}
	w.Header().Set("Content-Type", "application/json")

	info, err := io.ReadAll(r.Body)
	// check body

	if err != nil || len(info) == 0 {

		resp, _ := json.Marshal(createRespondBody(errBody))

		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", errBody,
		)
		return
	}
	// try unmarshal
	if err = json.Unmarshal(info, jsonURL); err != nil {
		resp, _ := json.Marshal(createRespondBody(errBody))

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(resp))
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", errBody,
		)
		return
	}
	// check valid url
	if !(utils.ValidateURL(jsonURL.URL)) {
		// change behaviour when requeust doesn't have the protocol 26-02-24
		jsonURL.URL = "http://" + jsonURL.URL
	}

	shortener = uh.store.GetShortener(jsonURL.URL)
	adr, _ := url.Parse(uh.baseURL)

	body = ""
	if len(adr.Host) == 0 {
		body = fmt.Sprintf("http://%s/%s", uh.baseURL, shortener)
	} else {
		body = fmt.Sprintf("%s/%s", uh.baseURL, shortener)
	}

	resp, err := json.Marshal(createRespondBody(body))
	if err != nil {
		resp, _ := json.Marshal(createRespondBody(errURL))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(resp))
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", errURL,
		)
		return
	}

	if len(uh.fileStore) > 0 {
		producer, err := NewProducer(uh.fileStore)
		if err != nil {
			sugar.Fatal(err)
		}
		defer producer.Close()
		if err := producer.WriteEvent(&Event{
			UUID:       uint(uh.store.Length()),
			ShortURL:   shortener,
			OriginlURL: jsonURL.URL,
		}); err != nil {
			sugar.Fatal(err)
		}
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(resp))

}

func createRespondBody(result string) JSONRespond {
	var respExtend = &JSONRespond{
		Result: result,
	}
	return *respExtend
}

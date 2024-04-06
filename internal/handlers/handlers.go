package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/closable/go-yandex-shortener/internal/utils"
	"go.uber.org/zap"
)

type Storager interface {
	GetShortener(userID int, txtURL string) (string, error)
	FindExistingKey(keyText string) (string, bool)
	Ping() bool
	PrepareStore()
	GetURLs(userID int) (map[string]string, error)
	SoftDeleteURLs(userID int, key ...string) error
}

type (
	URLHandler struct {
		store     Storager
		baseURL   string
		logger    zap.Logger
		maxLength int64
	}
	JSONRequest struct {
		URL string `json:"url"`
	}
	JSONRespond struct {
		Result string `json:"result"`
	}
	JSONBatch struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}
	JSONBatchRespond struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}
)

var (
	errBody    = "Error! Request body is empty!"
	errURL     = "Error! Check url it's should seems as like this 'http[s]://example.com'"
	emptyID    = "Error! id is empty!"
	notFoundID = "Error! id is not found or empty"
)

func New(st Storager, baseURL string, logger zap.Logger, maxLength int64) *URLHandler {
	st.PrepareStore()
	// load stored data
	// if len(fileStore) > 0 {
	// 	consumer, err := NewConsumer(fileStore)
	// 	if err != nil {
	// 		logger.Fatal("File not found")
	// 	}
	// 	defer consumer.file.Close()
	// 	loadDataFromFile(st, consumer.file)
	// }

	return &URLHandler{
		store:     st,
		baseURL:   baseURL,
		logger:    logger,
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
	statusSet := false
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

	w.Header().Set("Content-Type", "text/plain")
	userID, err := strconv.Atoi(r.FormValue("userID"))
	if err != nil {
		userID = 0
	}

	shortener, err = uh.store.GetShortener(userID, string(info))
	if err != nil {
		if err.Error() != "409" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errBody))
			sugar.Infoln(
				"uri", r.RequestURI,
				"method", r.Method,
				"description", errBody,
			)
			return
		}

		statusSet = true
		w.WriteHeader(http.StatusConflict)
	}
	body := makeShortenURL(shortener, uh.baseURL)
	if !statusSet {
		w.WriteHeader(http.StatusCreated)
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

	if ok && len(url) == 0 {
		w.WriteHeader(http.StatusGone)
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", "record was deleted",
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
	statusSet := false
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
	userID, err := strconv.Atoi(r.FormValue("userID"))
	if err != nil {
		userID = 0
	}
	shortener, err = uh.store.GetShortener(userID, jsonURL.URL)
	if err != nil {

		if err.Error() != "409" {
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
		statusSet = true
		w.WriteHeader(http.StatusConflict)
	}

	body = makeShortenURL(shortener, uh.baseURL)
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
	if !statusSet {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write([]byte(resp))

}

func (uh *URLHandler) UploadBatch(w http.ResponseWriter, r *http.Request) {
	sugar := *uh.logger.Sugar()
	var jsonData = &[]JSONBatch{}
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
	if err = json.Unmarshal(info, jsonData); err != nil {
		resp, _ := json.Marshal(createRespondBody(errBody))

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(resp))
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", err,
		)
		return
	}
	var bacthResp = []JSONBatchRespond{}
	var URL string
	userID, err := strconv.Atoi(r.FormValue("userID"))
	if err != nil {
		userID = 0
	}

	if len(*jsonData) > 0 {
		for _, v := range *jsonData {
			if !(utils.ValidateURL(v.OriginalURL)) {
				URL = "http://" + v.OriginalURL
			} else {
				URL = v.OriginalURL
			}

			shortener, _ := uh.store.GetShortener(userID, URL)
			body := makeShortenURL(shortener, uh.baseURL)

			item := &JSONBatchRespond{
				CorrelationID: v.CorrelationID,
				ShortURL:      body,
			}
			bacthResp = append(bacthResp, *item)
		}
	}

	resp, err := json.Marshal(bacthResp)
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

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(resp))

}

func makeShortenURL(URL string, baseURL string) string {
	adr, _ := url.Parse(baseURL)
	if len(adr.Host) == 0 {
		return fmt.Sprintf("http://%s/%s", baseURL, URL)
	} else {
		return fmt.Sprintf("%s/%s", baseURL, URL)
	}
}

func createRespondBody(result string) JSONRespond {
	var respExtend = &JSONRespond{
		Result: result,
	}
	return *respExtend
}

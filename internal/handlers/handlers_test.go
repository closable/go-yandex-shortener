package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/closable/go-yandex-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fileStore string = "/tmp/short-url-db.json"

var DSN string = ""

func TestGenerateShortener(t *testing.T) {
	// if len(DSN) == 0 {
	// 	cfg := config.LoadConfig()
	// 	DSN = cfg.DSN
	// }

	logger := NewLogger()
	// store, _ := storage.NewDBMS(DSN)
	// store, _ := storage.NewFile(fileStore)
	store, _ := storage.NewMemory()
	handler := New(store, "localhost:8080", logger, 1)
	store.Urls["test"] = "http://test.ru"

	type wants struct {
		method      string
		body        string
		contentType string
		statusCode  int
	}

	tests := []struct {
		name  string
		wants wants
	}{
		// TODO: Add test cases.
		{
			name: "Method POST",
			wants: wants{
				method:      "POST",
				body:        "https://yandex.ru",
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
			},
		},
		{
			name: "Method POST 409",
			wants: wants{
				method:      "POST",
				body:        "http://test.ru",
				contentType: "text/plain",
				statusCode:  http.StatusConflict,
			},
		},
		{
			name: "Method POST",
			wants: wants{
				method:      "POST",
				body:        "https://ffff.yandex.ru",
				contentType: "application/json",
				statusCode:  http.StatusCreated,
			},
		},
		{
			name: "Method POST fix bad request",
			wants: wants{
				method:      "POST",
				body:        "dddd.yandex.ru",
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
			},
		},
	}

	for _, tt := range tests {
		bodyReader := strings.NewReader(tt.wants.body)

		r := httptest.NewRequest(tt.wants.method, "/", bodyReader)
		w := httptest.NewRecorder()
		// вызовем хендлер как обычную функцию, без запуска самого сервера

		handler.GenerateShortener(w, r)
		assert.Equal(t, tt.wants.statusCode, w.Code, "Differents status codes")
		//assert.Equal(t, tt.wants.contentType, w.Header().Get("Content-Type"))

	}
}

func TestGetEndpointByShortener(t *testing.T) {
	// if len(DSN) == 0 {
	// 	cfg := config.LoadConfig()
	// 	DSN = cfg.DSN
	// }

	//store := storage.New()
	logger := NewLogger()
	// store, _ := storage.NewDBMS(DSN)
	// store, _ := storage.NewFile(fileStore)
	store, _ := storage.NewMemory()
	handler := New(store, "localhost:8080", logger, 1)

	type wants struct {
		method      string
		body        string
		contentType string
		statusCode  int
	}

	tests := []struct {
		name  string
		wants wants
	}{
		// TODO: Add test cases.
		{
			name: "Method GET",
			wants: wants{
				method:      "GET",
				body:        "https://yandex.ru",
				contentType: "text/plain",
				statusCode:  http.StatusTemporaryRedirect,
			},
		},
		{
			name: "Method GET bad request",
			wants: wants{
				method:      "GET",
				body:        "http://yandex.ru",
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		shortener := ""

		// create shortener prepare
		bodyReader := strings.NewReader(tt.wants.body)
		r := httptest.NewRequest("POST", "/", bodyReader)
		w := httptest.NewRecorder()

		handler.GenerateShortener(w, r)

		body := w.Body.String()
		path, err := url.Parse(body)
		require.Nil(t, err)
		shortener = path.Path[1:]

		if tt.wants.statusCode == http.StatusBadRequest {
			shortener = ""
		}

		// check result

		r = httptest.NewRequest(tt.wants.method, fmt.Sprintf("/%s", shortener), nil)
		w = httptest.NewRecorder()
		handler.GetEndpointByShortener(w, r)
		// вызовем хендлер как обычную функцию, без запуска самого сервера

		assert.Equal(t, tt.wants.statusCode, w.Code, "Differents status codes")
		//assert.Equal(t, tt.wants.contentType, w.Header().Get("Content-Type"))

	}
}

func TestGenerateJSONShortener(t *testing.T) {
	// if len(DSN) == 0 {
	// 	cfg := config.LoadConfig()
	// 	DSN = cfg.DSN
	// }

	logger := NewLogger()
	// store, _ := storage.NewDBMS(DSN)
	// store, _ := storage.NewFile(fileStore)
	store, _ := storage.NewMemory()
	handler := New(store, "localhost:8080", logger, 1)

	type wants struct {
		method      string
		body        string
		contentType string
		statusCode  int
	}

	tests := []struct {
		name  string
		wants wants
	}{
		// TODO: Add test cases.
		{
			name: "Method POST",
			wants: wants{
				method:      "POST",
				body:        "https://yandex.ru",
				contentType: "application/json",
				statusCode:  http.StatusCreated,
			},
		},
		// {
		// 	name: "Method POST",
		// 	wants: wants{
		// 		method:      "POST",
		// 		body:        "https://yandex.ru",
		// 		contentType: "application/json",
		// 		statusCode:  http.StatusConflict,
		// 	},
		// },
		{
			name: "Method POST wrong content-type",
			wants: wants{
				method:      "POST",
				body:        "https://zzzz.yandex.ru",
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
			},
		},
		{
			name: "Method POST after bad request",
			wants: wants{
				method:      "POST",
				body:        "ffff.yandex.ru",
				contentType: "application/json",
				statusCode:  http.StatusCreated,
			},
		},
	}
	for _, tt := range tests {

		var jsonURL = &JSONRequest{
			URL: tt.wants.body,
		}
		body, _ := json.Marshal(jsonURL)
		bodyReader := bytes.NewReader([]byte(body))

		r := httptest.NewRequest(tt.wants.method, "/api/shorten", bodyReader)
		w := httptest.NewRecorder()
		// вызовем хендлер как обычную функцию, без запуска самого сервера

		handler.GenerateJSONShortener(w, r)

		assert.Equal(t, tt.wants.statusCode, w.Code, "Differents status codes")

		if tt.wants.contentType != "application/json" {
			assert.NotEqual(t, tt.wants.contentType, w.Header().Get("Content-Type"), "Wrong content-type")
		} else {
			assert.Equal(t, tt.wants.contentType, w.Header().Get("Content-Type"), "Wrong content-type")
		}
	}
}

func TestCompressor(t *testing.T) {
	// if len(DSN) == 0 {
	// 	cfg := config.LoadConfig()
	// 	DSN = cfg.DSN
	// }

	logger := NewLogger()
	// store, _ := storage.NewDBMS(DSN)
	// store, _ := storage.NewFile(fileStore)
	store, _ := storage.NewMemory()
	handler := New(store, "localhost:8080", logger, 1)

	ts := httptest.NewServer(handler.InitRouter())
	defer ts.Close()

	tests := []struct {
		name              string
		path              string
		body              string
		expectedEncoding  string
		acceptedEncodings string
	}{
		{
			name:              "equal encodings",
			path:              "/",
			body:              "http://mmm.yandex.ru",
			acceptedEncodings: "gzip",
			expectedEncoding:  "gzip",
		},
		{
			name:              "equal encodings JSON",
			path:              "/api/shorten",
			body:              "{\"url\": \"http://ffff.yandex.ru\"}",
			acceptedEncodings: "gzip",
			expectedEncoding:  "gzip",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			body := strings.NewReader(tc.body)

			r, _ := http.NewRequest("POST", ts.URL+tc.path, body)

			r.Header.Set("Accept-Encoding", "gzip")

			resp, err := http.DefaultClient.Do(r)
			require.NoError(t, err)

			defer resp.Body.Close()

			// b, err := io.ReadAll(resp.Body)
			// require.NoError(t, err)
			// fmt.Println(string(b))

			require.Equal(t, tc.expectedEncoding, resp.Header.Get("Accept-Encoding"))

		})
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			body := strings.NewReader(tc.body)

			r, _ := http.NewRequest("POST", ts.URL+tc.path, body)
			r.Header.Set("Accept-Encoding", "gzip")

			resp, err := http.DefaultClient.Do(r)
			require.NoError(t, err)
			defer resp.Body.Close()

			zr, err := gzip.NewReader(resp.Body)
			require.NoError(t, err)

			b, err := io.ReadAll(zr)
			require.NoError(t, err)

			// require.Equal(t, resp.StatusCode, 201)
			require.Equal(t, tc.expectedEncoding, resp.Header.Get("Accept-Encoding"))
			fmt.Println(string(b))

		})
	}

}

func TestUploadBatch(t *testing.T) {
	// if len(DSN) == 0 {
	// 	cfg := config.LoadConfig()
	// 	DSN = cfg.DSN
	// }

	logger := NewLogger()
	// store, _ := storage.NewDBMS(DSN)
	// store, _ := storage.NewFile(fileStore)
	store, _ := storage.NewMemory()
	handler := New(store, "localhost:8080", logger, 1)

	type wants struct {
		method      string
		body        string
		contentType string
		statusCode  int
		resultLen   int
	}

	tests := []struct {
		name  string
		wants wants
	}{
		// TODO: Add test cases.
		{
			name: "Method POST",
			wants: wants{
				method: "POST",
				body: `[
					{
						"correlation_id": "cudLI",
						"original_url": "http://wSKr/yandex.ru/uuVGbXV"
					},
					{
						"correlation_id": "mUOvN",
						"original_url": "http://Rktt/yandex.ru/xxLOsfQ"
					},
					{
						"correlation_id": "SriwT",
						"original_url": "http://XVap/yandex.ru/QWdMOMm"
					}]`,
				contentType: "application/json",
				statusCode:  http.StatusCreated,
				resultLen:   3,
			},
		},
	}
	for _, tt := range tests {

		bodyReader := bytes.NewReader([]byte(tt.wants.body))

		r := httptest.NewRequest(tt.wants.method, "/api/shorten/batch", bodyReader)
		w := httptest.NewRecorder()
		// вызовем хендлер как обычную функцию, без запуска самого сервера

		handler.UploadBatch(w, r)
		assert.Equal(t, tt.wants.resultLen, store.Length(), "The result of the operation is wrong")
		assert.Equal(t, tt.wants.statusCode, w.Code, "Differents status codes")

		if tt.wants.contentType != "application/json" {
			assert.NotEqual(t, tt.wants.contentType, w.Header().Get("Content-Type"), "Wrong content-type")
		} else {
			assert.Equal(t, tt.wants.contentType, w.Header().Get("Content-Type"), "Wrong content-type")
		}
	}
}

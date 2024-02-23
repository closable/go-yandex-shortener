package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/closable/go-yandex-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateShortener(t *testing.T) {
	store := storage.New()
	logger := NewLogger()
	handler := New(store, "localhost:8080", logger)

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
			name: "Method POST",
			wants: wants{
				method:      "POST",
				body:        "https://yandex.ru",
				contentType: "application/json",
				statusCode:  http.StatusCreated,
			},
		},
		{
			name: "Method POST bad request",
			wants: wants{
				method:      "POST",
				body:        "yandex.ru",
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
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
	store := storage.New()
	logger := NewLogger()
	handler := New(store, "localhost:8080", logger)

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
	store := storage.New()
	logger := NewLogger()
	handler := New(store, "localhost:8080", logger)

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
		{
			name: "Method POST wrong content-type",
			wants: wants{
				method:      "POST",
				body:        "https://yandex.ru",
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
			},
		},
		{
			name: "Method POST bad request",
			wants: wants{
				method:      "POST",
				body:        "yandex.ru",
				contentType: "application/json",
				statusCode:  http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {

		var jsonUrl = &JSONRequest{
			URL: tt.wants.body,
		}
		body, _ := json.Marshal(jsonUrl)
		bodyReader := bytes.NewReader([]byte(body))

		r := httptest.NewRequest(tt.wants.method, "/shorten", bodyReader)
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

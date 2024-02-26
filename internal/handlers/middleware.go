package handlers

import (
	"compress/gzip"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

func (uh *URLHandler) Logger(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		sugar := *uh.logger.Sugar()
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter
		duration := time.Since(start)

		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status, // получаем перехваченный код статуса ответа
			"duration", duration,
			"size", responseData.size, // получаем перехваченный размер ответа
		)
	}

	return http.HandlerFunc(logFn)
}

func (uh *URLHandler) Compressor(h http.Handler) http.Handler {
	zipFn := func(w http.ResponseWriter, r *http.Request) {
		sugar := *uh.logger.Sugar()
		contLength, _ := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") || contLength < uh.maxLength {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			sugar.Infoln(
				"uri", r.RequestURI,
				"method", r.Method,
				"content", "default",
			)

			h.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"content", "gzip",
		)

		h.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)

	}

	return http.HandlerFunc(zipFn)
}

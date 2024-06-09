// Package handlers реализует функцию для упровления handlers
package handlers

import (
	"compress/gzip"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TokenEXP Константа для создания токена
const TokenEXP = time.Hour * 3

// SecretKEY Константа для создания токена
const SecretKEY = "*HelloWorld*"

// описание структур данных
type (
	//gzipWriter struct Структура для работы compressor
	gzipWriter struct {
		http.ResponseWriter
		Writer io.Writer
	}
	// Claims Структура для работы autheticator
	Claims struct {
		jwt.RegisteredClaims
		UserID int
	}
)

// Write вспомогательная функция для реализации сжатия информации
func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

// Logger middleware для ведения логгирования запросов
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

// Compressor middleware для сжатия запроса
func (uh *URLHandler) Compressor(h http.Handler) http.Handler {
	zipFn := func(w http.ResponseWriter, r *http.Request) {
		sugar := *uh.logger.Sugar()
		//contLength, _ := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") /*|| contLength < uh.maxLength*/ {
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

		w.Header().Set("Accept-Encoding", "gzip")
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

// Auth аутентификация пользователя
func (uh *URLHandler) Auth(h http.Handler) http.Handler {
	auth := func(w http.ResponseWriter, r *http.Request) {
		sugar := *uh.logger.Sugar()
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/api/user/urls" || r.URL.Path == "/" {
			var userID int
			token, errCookie := r.Cookie("Authorization")
			headerAuth := r.Header.Get("Authorization")
			fmt.Printf("-1 %s 2- %s 3- %s ", token, errCookie, headerAuth)
			// if errCookie has err && headerAuth empty
			if errCookie != nil && len(headerAuth) == 0 {
				tokenString, err := BuildJWTString()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(" "))
					sugar.Infoln(
						"uri", r.RequestURI,
						"method", r.Method,
						"description", err,
					)
					h.ServeHTTP(w, r)
					return
				}

				cookie := http.Cookie{
					Name:    "Authorization",
					Expires: time.Now().Add(TokenEXP),
					Value:   tokenString,
				}
				http.SetCookie(w, &cookie)
				w.Header().Add("Authorization", tokenString)
				userID = GetUserID(tokenString)
				fmt.Printf("user get from empty cookies %d\n", userID)
			}

			if len(token.String()) > 0 {
				userID = GetUserID(token.Value)
				w.Header().Add("Authorization", token.Value)
				fmt.Printf("user get from existing cookies %d\n", userID)
			}

			if len(headerAuth) > 0 && userID == 0 {
				userID = GetUserID(headerAuth)
				w.Header().Add("Authorization", headerAuth)
				fmt.Printf("user get from existing header %d\n", userID)
			}

			//userID = 0
			values := url.Values{}
			values.Add("userID", fmt.Sprintf("%d", userID))
			r.PostForm = values

			if userID == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(" "))
				sugar.Infoln(
					"uri", r.RequestURI,
					"method", r.Method,
					"description", "User unauthorized",
				)
				h.ServeHTTP(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(auth)
}

// BuildJWTString посроение строки токена аутетификации
func BuildJWTString() (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenEXP)),
		},
		// собственное утверждение
		UserID: rand.Intn(1000),
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(SecretKEY))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

// GetUserID вспомогательная функция для определения UserID из токена
func GetUserID(tokenString string) int {
	// создаём экземпляр структуры с утверждениями
	claims := &Claims{}
	// парсим из строки токена tokenString в структуру claims
	jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKEY), nil
	})

	// возвращаем ID пользователя в читаемом виде
	return claims.UserID
}

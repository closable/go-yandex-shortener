package handlers

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const TokenEXP = time.Hour * 3
const SecretKEY = "*HelloWorld*"

type (
	gzipWriter struct {
		http.ResponseWriter
		Writer io.Writer
	}
	Claims struct {
		jwt.RegisteredClaims
		UserID int
	}
)

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

func (uh *URLHandler) Auth(h http.Handler) http.Handler {
	auth := func(w http.ResponseWriter, r *http.Request) {
		sugar := *uh.logger.Sugar()

		if r.URL.Path == "/api/user/urls" {
			var userID int
			token, err := r.Cookie("Authorization")
			if err != nil {
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
				userID = GetUserID(tokenString)
				fmt.Printf("user get from empty cookies %d", userID)
			}

			if len(token.String()) > 0 {
				userID = GetUserID(token.Value)
				fmt.Printf("user get from existing cookies %d", userID)
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

func BuildJWTString() (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenEXP)),
		},
		// собственное утверждение
		UserID: 13,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(SecretKEY))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

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

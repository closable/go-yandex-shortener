package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type AuthItemURL struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

func (uh *URLHandler) GetUrls(w http.ResponseWriter, r *http.Request) {
	sugar := *uh.logger.Sugar()

	var userId int
	token, err := r.Cookie("Authorization")

	if err != nil {
		tokenString, err := BuildJWTString()
		if err != nil {
			log.Fatal(err)
		}

		cookie := http.Cookie{
			Name:    "Authorization",
			Expires: time.Now().Add(TokenEXP),
			Value:   tokenString,
		}
		http.SetCookie(w, &cookie)
		userId = GetUserID(tokenString)
	}

	if len(token.String()) > 0 {
		userId = GetUserID(token.Value)
	}

	w.Header().Set("Content-Type", "application/json")

	var body []AuthItemURL
	res, err := uh.store.GetURLs(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errBody))
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", err,
		)
		return
	}

	if len(res) == 0 {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(" "))
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", fmt.Sprintf("recs for user %d not found", userId),
		)
		return
	}

	var shorten string

	for key, value := range res {
		shorten = makeShortenURL(key, uh.baseURL)
		body = append(body, makeAuthRequestRow(shorten, value))
	}
	resp, err := json.Marshal(body)
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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp))
}

func makeAuthRequestRow(key, url string) AuthItemURL {
	var res = &AuthItemURL{
		ShortURL:    key,
		OriginalURL: url,
	}

	return *res
}

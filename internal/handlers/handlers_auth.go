package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type (
	// Структура для работы JSON
	AuthItemURL struct {
		OriginalURL string `json:"original_url"`
		ShortURL    string `json:"short_url"`
	}
)

// GetUrls функция для получения urls пользователя
func (uh *URLHandler) GetUrls(w http.ResponseWriter, r *http.Request) {
	sugar := *uh.logger.Sugar()
	userID, _ := strconv.Atoi(r.FormValue("userID"))

	w.Header().Set("Content-Type", "application/json")

	var body []AuthItemURL
	res, err := uh.store.GetURLs(int(userID))
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
		//w.WriteHeader(http.StatusNoContent)
		w.WriteHeader(http.StatusUnauthorized) // ну как так?, не согласен
		w.Write([]byte(" "))
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", fmt.Sprintf("recs for user %d not found", userID),
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

// DelUrls функция для удаления указанных значений
func (uh *URLHandler) DelUrls(w http.ResponseWriter, r *http.Request) {
	sugar := *uh.logger.Sugar()
	userID, _ := strconv.Atoi(r.FormValue("userID"))

	info, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	var kyesForDelete = make([]string, 0)

	if err = json.Unmarshal(info, &kyesForDelete); err != nil {
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(" "))
		return
	}

	if len(kyesForDelete) == 0 {
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", "Nothing to delete",
		)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(" "))
		return
	}

	if len(kyesForDelete) > 0 {
		err := uh.store.SoftDeleteURLs(userID, kyesForDelete...)
		if err != nil {
			sugar.Infoln(
				"uri", r.RequestURI,
				"method", r.Method,
				"description", err,
			)
		}
	}

	sugar.Infoln(
		"uri", r.RequestURI,
		"method", r.Method,
		"description", "selected records where deleted successfully",
	)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(" "))
}

// makeAuthRequestRow функция помошник для составления тела запроса
func makeAuthRequestRow(key, url string) AuthItemURL {
	var res = &AuthItemURL{
		ShortURL:    key,
		OriginalURL: url,
	}

	return *res
}

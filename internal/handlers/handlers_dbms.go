package handlers

import (
	"encoding/json"
	"net/http"
)

var (
	connOk  = "The connection is still alive"
	connErr = "The connection was lost"
)

func (uh URLHandler) CloseBaseActivity(w http.ResponseWriter, r *http.Request) {
	uh.conn.Close(uh.ctx)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(createRespondBody("Connection was closed"))
	w.Write([]byte(resp))

}

func (uh URLHandler) CheckBaseActivity(w http.ResponseWriter, r *http.Request) {
	sugar := *uh.logger.Sugar()
	w.Header().Set("Content-Type", "application/json")

	err := uh.conn.Ping(uh.ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(createRespondBody(connErr))
		w.Write([]byte(resp))
		if err != nil {
			sugar.Infoln(
				"uri", r.RequestURI,
				"method", r.Method,
				"description", connErr,
			)
		}
		return
	}

	resp, err := json.Marshal(createRespondBody(connOk))
	if err != nil {
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", connOk,
		)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp))

}

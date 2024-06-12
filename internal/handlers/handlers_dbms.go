package handlers

import (
	"encoding/json"
	"net/http"
)

var (
	connOk  = "The connection is still alive"
	connErr = "The connection was lost"
)

// CheckBaseActivity haelth check
func (uh URLHandler) CheckBaseActivity(w http.ResponseWriter, r *http.Request) {
	sugar := *uh.logger.Sugar()
	w.Header().Set("Content-Type", "application/json")
	ok := uh.store.Ping()

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(createRespondBody(connErr))
		w.Write([]byte(resp))
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", connErr,
		)
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

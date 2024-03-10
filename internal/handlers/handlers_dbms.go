package handlers

import (
	"encoding/json"
	"net/http"
)

var (
	connOk  = "The connection is still alive"
	connErr = "The connection was lost"
)

func (uh URLHandler) CheckBaseActivity(w http.ResponseWriter, r *http.Request) {
	sugar := *uh.logger.Sugar()
	w.Header().Set("Content-Type", "application/json")
	conn, err := uh.dbms.GetConn()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(createRespondBody(connErr))
		w.Write([]byte(resp))
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", err,
		)
		return
	}
	defer conn.Close()
	err = conn.PingContext(uh.dbms.CTX)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(createRespondBody(connErr))
		w.Write([]byte(resp))
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"description", err, //connErr,
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

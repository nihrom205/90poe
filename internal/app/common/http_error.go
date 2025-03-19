package common

import (
	"encoding/json"
	"net/http"
)

func NotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w)
}

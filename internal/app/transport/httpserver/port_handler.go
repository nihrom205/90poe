package httpserver

import (
	"github.com/gorilla/mux"
	"github.com/nihrom205/90poe/internal/app/common"
	"net/http"
)

func (h HttpServer) Processing(w http.ResponseWriter, r *http.Request) {
	h.portService.ProcessingJson(r.Context(), r.Body)
	common.RespondOK("", w)
}

func (h HttpServer) GetPortByKey(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	port, err := h.portService.GetPortByKey(key)
	if err != nil {
		common.NotFound(w)
	}
	common.RespondOK(port, w)
}

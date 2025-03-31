package httpserver

import (
	"github.com/gorilla/mux"
	"github.com/nihrom205/90poe/internal/app/common"
	"net/http"
)

func (h HttpServer) LoadPorts(w http.ResponseWriter, r *http.Request) {
	h.portService.UploadPorts(r.Context(), r.Body)
	common.RespondOK("ok: true", w)
}

func (h HttpServer) GetPort(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	port, err := h.portService.GetPort(r.Context(), key)
	if err != nil {
		common.NotFound(w)
	}
	common.RespondOK(port, w)
}

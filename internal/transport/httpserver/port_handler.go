package httpserver

import "net/http"

func (h HttpServer) MyTestHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}

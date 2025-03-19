package httpserver

import (
	"github.com/nihrom205/90poe/internal/app/service"
)

type HttpServer struct {
	portService IPortService
}

func NewHttpServer(portService service.PortService) *HttpServer {
	return &HttpServer{
		portService: portService,
	}
}

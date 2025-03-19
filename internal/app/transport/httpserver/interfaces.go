package httpserver

import (
	"context"
	"github.com/nihrom205/90poe/internal/app/service"
	"io"
)

type IPortService interface {
	ProcessingJson(ctx context.Context, data io.ReadCloser)
	GetPortByKey(key string) (*service.Port, error)
}

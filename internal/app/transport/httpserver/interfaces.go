package httpserver

import (
	"context"
	"io"

	"github.com/nihrom205/90poe/internal/app/service"
)

type IPortService interface {
	ProcessingJson(ctx context.Context, data io.ReadCloser)
	GetPortByKey(key string) (*service.Port, error)
}

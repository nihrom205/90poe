package httpserver

import (
	"context"
	"io"

	"github.com/nihrom205/90poe/internal/app/domain"
)

type IPortService interface {
	UploadPorts(ctx context.Context, data io.ReadCloser)
	GetPort(ctx context.Context, key string) (*domain.NewPortData, error)
}

package httpserver

import (
	"context"
	"io"

	"github.com/nihrom205/90poe/internal/app/domain"
)

type IPortService interface {
	ProcessingJson(ctx context.Context, data io.ReadCloser)
	GetPortByKey(ctx context.Context, key string) (*domain.NewPortData, error)
}

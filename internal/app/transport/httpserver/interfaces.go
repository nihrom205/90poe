package httpserver

import (
	"context"
	"github.com/nihrom205/90poe/internal/app/domain"
	"io"
)

type IPortService interface {
	ProcessingJson(ctx context.Context, data io.ReadCloser)
	GetPortByKey(ctx context.Context, key string) (*domain.NewPortData, error)
}

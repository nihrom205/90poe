package service

import (
	"context"
	"github.com/nihrom205/90poe/internal/app/repository"
)

type IPortRepository interface {
	CreatePort(ctx context.Context, port *repository.Port) (*repository.Port, error)
	UpdateLocation(ctx context.Context, port *repository.Port) (*repository.Port, error)
	GetPortByKey(ctx context.Context, key string) (*repository.Port, error)
	GetAllPorts(ctx context.Context) ([]repository.Port, error)
}

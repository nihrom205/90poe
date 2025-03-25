package service

import (
	"context"
	"github.com/nihrom205/90poe/internal/app/repository"
	"github.com/stretchr/testify/mock"
)

type MockPortRepository struct {
	mock.Mock
}

func (m *MockPortRepository) CreatePort(ctx context.Context, port *repository.Port) (*repository.Port, error) {
	args := m.Called(port)
	return args.Get(0).(*repository.Port), args.Error(1)
}

func (m *MockPortRepository) UpdateLocation(ctx context.Context, port *repository.Port) (*repository.Port, error) {
	args := m.Called(port)
	return args.Get(0).(*repository.Port), args.Error(1)
}

func (m *MockPortRepository) GetPortByKey(ctx context.Context, key string) (*repository.Port, error) {
	args := m.Called(key)
	return args.Get(0).(*repository.Port), args.Error(1)
}

func (m *MockPortRepository) GetAllPorts(ctx context.Context) ([]repository.Port, error) {
	args := m.Called()
	return args.Get(0).([]repository.Port), args.Error(1)
}

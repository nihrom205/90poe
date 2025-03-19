package service

import "github.com/nihrom205/90poe/internal/app/repository"

type IPortRepository interface {
	CreatePort(port *repository.Port) (*repository.Port, error)
	UpdateLocation(port *repository.Port) (*repository.Port, error)
	GetPortByKey(key string) (*repository.Port, error)
	GetAllPorts() ([]repository.Port, error)
}

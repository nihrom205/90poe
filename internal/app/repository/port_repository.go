package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/nihrom205/90poe/internal/pkg/pg"

	"github.com/nihrom205/90poe/internal/app/domain"
	"gorm.io/gorm"
)

type PortRepository struct {
	Db *pg.Db
}

func NewPortRepository(db *pg.Db) PortRepository {
	return PortRepository{Db: db}
}

func (repo PortRepository) CreatePort(ctx context.Context, port *Port) (*Port, error) {
	result := repo.Db.Create(port)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create port: %w", result.Error)
	}
	return port, nil
}

func (repo PortRepository) UpdateLocation(ctx context.Context, port *Port) error {
	fmt.Println("UpdateLocation start= ", port)
	result := repo.Db.Save(port)
	if result.Error != nil {
		return fmt.Errorf("failed to update location: %w", result.Error)
	}
	fmt.Println("UpdateLocation end= ", port)
	return nil
}

func (repo PortRepository) GetPort(ctx context.Context, key string) (*Port, error) {
	var port Port
	result := repo.Db.First(&port, "key = ?", key)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get port by key: %w", result.Error)
	}
	return &port, nil
}

func (repo PortRepository) GetAllPorts(ctx context.Context) ([]Port, error) {
	var ports []Port
	err := repo.Db.Find(&ports).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return ports, fmt.Errorf("failed to get all ports: %w", err)
	}
	return ports, nil
}

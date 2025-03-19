package repository

import (
	"github.com/nihrom205/90poe/internal/pkg"
)

type PortRepository struct {
	Db *pkg.Db
}

func NewPortRepository(db *pkg.Db) *PortRepository {
	return &PortRepository{Db: db}
}

func (repo PortRepository) CreatePort(port *Port) (*Port, error) {
	result := repo.Db.Create(port)
	if result.Error != nil {
		return nil, result.Error
	}
	return port, nil
}

func (repo PortRepository) UpdateLocation(port *Port) (*Port, error) {
	result := repo.Db.Save(port)
	if result.Error != nil {
		return nil, result.Error
	}
	return port, nil
}

func (repo PortRepository) GetPortByKey(key string) (*Port, error) {
	var port Port
	result := repo.Db.First(&port, "key = ?", key)
	if result.Error != nil {
		return nil, result.Error
	}
	return &port, nil
}

//func (repo PortRepository) GetPortById(id int) (Port, error) {
//	var port Port
//	result := repo.Db.Find(&port, "id = ?", id)
//	if result.Error != nil {
//		fmt.Println(result.Error)
//	}
//	return port, nil
//}

func (repo PortRepository) GetAllPorts() ([]Port, error) {
	var ports []Port
	if err := repo.Db.Find(&ports).Error; err != nil {
		return nil, err
	}
	return ports, nil
}

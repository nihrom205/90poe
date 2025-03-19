package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nihrom205/90poe/internal/app/repository"
	"gorm.io/gorm"
	"io"
	"log"
	"sync"
	"time"
)

type PortService struct {
	repo IPortRepository
}

func NewPortService(repo IPortRepository) PortService {
	return PortService{
		repo: repo,
	}
}

func (s PortService) ProcessingJson(ctx context.Context, data io.ReadCloser) {
	chLocation := make(chan keyAndLocation, 1)
	var wg sync.WaitGroup
	//g := new(errgroup.Group)

	wg.Add(1)
	go func() {
		defer wg.Done()
		SavePort(ctx, s.repo, chLocation)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		Parse(ctx, chLocation, data)
	}()

	wg.Wait()
}

func Parse(ctx context.Context, ch chan<- keyAndLocation, data io.ReadCloser) {
	defer close(ch)
	decoder := json.NewDecoder(data)

	if _, err := decoder.Token(); err != nil {
		log.Fatalf("PortService - Parse: Ошибка при чтении открывающего токена: %v", err)
	}

	for decoder.More() {
		select {
		case <-ctx.Done():
			return
		default:
			key, err := decoder.Token()
			if err != nil {
				fmt.Println("PortService - Parse: Ошибка при чтении ключа:", err)
				continue
			}

			var location Location
			if err = decoder.Decode(&location); err != nil {
				fmt.Println("PortService - Parse: Ошибка при декодировании объекта:", err)
				continue
			}

			ch <- keyAndLocation{
				Key: key.(string),
				Loc: location,
			}
		}
	}

	if _, err := decoder.Token(); err != nil {
		fmt.Println("PortService - Parse: Ошибка при чтении токена:", err)
		return
	}
}

func SavePort(ctx context.Context, repo IPortRepository, key <-chan keyAndLocation) {
	for {
		select {
		case <-ctx.Done():
			return
		case keyLoc, ok := <-key:
			if !ok {
				log.Println("PortService - SavePort: канал закрыт, завершение выполнения")
				return
			}
			port, isValid := mapperToDB(keyLoc)
			if !isValid {
				log.Println("PortService - SavePort: Ошибка валидации")
				continue
			}

			portDb, err := repo.GetPortByKey(port.Key)
			if portDb == nil && errors.Is(err, gorm.ErrRecordNotFound) {
				_, err = repo.CreatePort(port)
				if err != nil {
					log.Printf("PortService - SavePort: не удалось сохранить: %v", err)
				}
				continue
			}

			port.ID = portDb.ID
			port.CreatedAt = portDb.CreatedAt
			port.UpdatedAt = time.Now()
			_, err = repo.UpdateLocation(port)
			if err != nil {
				log.Printf("PortService - SavePort: Ошибка при обновлении %v", err)
			}
		}
	}
}

func (s PortService) GetPortByKey(key string) (*Port, error) {
	portDb, err := s.repo.GetPortByKey(key)
	if err != nil {
		return nil, err
	}
	port := mapperPort(portDb)
	return port, nil

}

func mapperToDB(keyLoc keyAndLocation) (*repository.Port, bool) {
	if keyLoc.Key == "" {
		return nil, false
	}

	return &repository.Port{
		Key:         keyLoc.Key,
		Name:        keyLoc.Loc.Name,
		City:        keyLoc.Loc.City,
		Country:     keyLoc.Loc.Country,
		Alias:       keyLoc.Loc.Alias,
		Regions:     keyLoc.Loc.Regions,
		Coordinates: keyLoc.Loc.Coordinates,
		Province:    keyLoc.Loc.Province,
		Timezone:    keyLoc.Loc.Timezone,
		Unlocs:      keyLoc.Loc.Unlocs,
		Code:        keyLoc.Loc.Code,
	}, true
}

func mapperPort(port *repository.Port) *Port {
	return &Port{
		Key:         port.Key,
		Name:        port.Name,
		City:        port.City,
		Country:     port.Country,
		Alias:       port.Alias,
		Regions:     port.Regions,
		Coordinates: port.Coordinates,
		Province:    port.Province,
		Timezone:    port.Timezone,
		Unlocs:      port.Unlocs,
		Code:        port.Code,
	}
}

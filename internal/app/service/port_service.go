package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"sync"
	"time"

	"github.com/nihrom205/90poe/internal/app/repository"
	"gorm.io/gorm"
)

type PortService struct {
	repo   IPortRepository
	logger *zerolog.Logger
}

func NewPortService(repo IPortRepository, logger *zerolog.Logger) PortService {
	return PortService{
		repo:   repo,
		logger: logger,
	}
}

func (s PortService) ProcessingJson(ctx context.Context, data io.ReadCloser) {
	chLocation := make(chan keyAndLocation, 1)
	var wg sync.WaitGroup
	//g := new(errgroup.Group)

	wg.Add(1)
	go func() {
		defer wg.Done()
		SavePort(ctx, s.logger, s.repo, chLocation)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		Parse(ctx, s.logger, chLocation, data)
	}()

	wg.Wait()
}

func Parse(ctx context.Context, logger *zerolog.Logger, ch chan<- keyAndLocation, data io.ReadCloser) {
	defer close(ch)
	decoder := json.NewDecoder(data)

	if _, err := decoder.Token(); err != nil {
		logger.Error().Msgf("PortService - Parse: Error reading open token: %v", err)
		return
	}

	for decoder.More() {
		select {
		case <-ctx.Done():
			return
		default:
			keyToken, err := decoder.Token()
			if err != nil {
				logger.Error().Msgf("PortService - Parse: Error read for keyToken: %v", err)
				continue
			}

			strKey, ok := keyToken.(string)
			if !ok {
				logger.Error().Msgf("PortService - Parse: Error converting keyToken to string: %v", err)
				continue
			}

			var location Location
			if err = decoder.Decode(&location); err != nil {
				logger.Error().Msgf("PortService - Parse: Error converting keyToken to string: %v", err)
				continue
			}

			ch <- keyAndLocation{
				Key: strKey,
				Loc: location,
			}
		}
	}

	if _, err := decoder.Token(); err != nil {
		logger.Error().Msgf("PortService - Parse: Error reading closing token: %v", err)
		return
	}
}

func SavePort(ctx context.Context, logger *zerolog.Logger, repo IPortRepository, key <-chan keyAndLocation) {
	for {
		select {
		case <-ctx.Done():
			return
		case keyLoc, ok := <-key:
			if !ok {
				logger.Info().Msgf("PortService - SavePort: Completion of execution")
				return
			}
			port, isValid := mapperToDB(keyLoc)
			if !isValid {
				logger.Error().Msgf("PortService - SavePort: Error validation port")
				continue
			}

			portDb, err := repo.GetPortByKey(port.Key)
			if portDb == nil && errors.Is(err, gorm.ErrRecordNotFound) {
				_, err = repo.CreatePort(port)
				if err != nil {
					logger.Error().Msgf("PortService - SavePort: Error creating port: %v", err)
				}
				continue
			}

			port.ID = portDb.ID
			port.CreatedAt = portDb.CreatedAt
			port.UpdatedAt = time.Now()
			_, err = repo.UpdateLocation(port)
			if err != nil {
				logger.Error().Msgf("PortService - SavePort: Error updating location: %v", err)
			}
		}
	}
}

func (s PortService) GetPortByKey(key string) (*Port, error) {
	portDb, err := s.repo.GetPortByKey(key)
	if err != nil {
		return nil, fmt.Errorf("PortService - GetPortByKey: %w", err)
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

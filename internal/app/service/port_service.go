package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/nihrom205/90poe/internal/app/domain"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/nihrom205/90poe/internal/app/repository"
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

func (s PortService) UploadPorts(ctx context.Context, data io.ReadCloser) {
	chLocation := make(chan keyAndLocation, 1)
	g := errgroup.Group{}

	g.Go(func() error {
		return SavePort(ctx, s.logger, s.repo, chLocation)
	})

	g.Go(func() error {
		return Parse(ctx, s.logger, chLocation, data)
	})

	if err := g.Wait(); err != nil {
		s.logger.Error().Err(err).Msg("PortService - UploadPorts: Error processing json")
	}
}

func Parse(ctx context.Context, logger *zerolog.Logger, ch chan<- keyAndLocation, data io.ReadCloser) error {
	defer close(ch)
	decoder := json.NewDecoder(data)

	if _, err := decoder.Token(); err != nil {
		logger.Error().Err(err).Msg("PortService - Parse: Error reading open token")
		return fmt.Errorf("PortService - Parse: Error reading open token: %w", err)
	}

	for decoder.More() {
		select {
		case <-ctx.Done():
			return nil
		default:
			keyToken, err := decoder.Token()
			if err != nil {
				logger.Error().Err(err).Msg("PortService - Parse: Error read for keyToken")
				continue
			}

			strKey, ok := keyToken.(string)
			if !ok {
				logger.Error().Err(err).Msg("PortService - Parse: Error converting keyToken to string")
				continue
			}

			var location Location
			if err = decoder.Decode(&location); err != nil {
				logger.Error().Err(err).Msg("PortService - Parse: Error converting keyToken to string")
				continue
			}

			ch <- keyAndLocation{
				Key: strKey,
				Loc: location,
			}
		}
	}

	if _, err := decoder.Token(); err != nil {
		logger.Error().Err(err).Msg("PortService - Parse: Error reading closing token")
		return fmt.Errorf("PortService - Parse: Error reading closing token: %w", err)
	}
	return nil
}

func SavePort(ctx context.Context, logger *zerolog.Logger, repo IPortRepository, key <-chan keyAndLocation) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case keyLoc, ok := <-key:
			if !ok {
				logger.Info().Msg("PortService - SavePort: Completion of execution")
				return nil
			}
			port, isValid := mapperToDB(keyLoc)
			if !isValid {
				logger.Error().Msg("PortService - SavePort: Error validation port")
				continue
			}

			portDb, err := repo.GetPort(ctx, port.Key)
			if portDb == nil && err != nil {
				_, err = repo.CreatePort(ctx, port)
				if err != nil {
					logger.Error().Err(err).Msg("PortService - SavePort: Error creating port")
				}
				continue
			}

			port.ID = portDb.ID
			port.CreatedAt = portDb.CreatedAt
			port.UpdatedAt = time.Now()
			err = repo.UpdateLocation(ctx, port)
			if err != nil {
				logger.Error().Err(err).Msg("PortService - SavePort: Error updating location")
			}
		}
	}
}

func (s PortService) GetPort(ctx context.Context, key string) (*domain.NewPortData, error) {
	portDb, err := s.repo.GetPort(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("PortService - GetPort: %w", err)
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

func mapperPort(port *repository.Port) *domain.NewPortData {
	return &domain.NewPortData{
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

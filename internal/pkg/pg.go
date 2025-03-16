package pkg

import (
	"errors"
	"fmt"
	"github.com/nihrom205/90poe/internal/app/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

type Db struct {
	*gorm.DB
}

func NewDb(config *config.Config) (*Db, error) {
	var dsn = config.DSN
	//var dsn = "file::memory:?cache=shared"
	if dsn == "" {
		return nil, errors.New("нет DSN")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("не удалось получить объект sql.DB: %v", err)
	}

	sqlDb.SetMaxOpenConns(10)
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetConnMaxLifetime(1 * time.Minute)

	if err = sqlDb.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping failed: %w", err)
	}
	return &Db{DB: db}, nil
}

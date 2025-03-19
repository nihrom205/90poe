package pkg

import (
	"errors"
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"log"
	"time"
)

type Db struct {
	*gorm.DB
}

func NewDb(dsn string) (*Db, error) {
	if dsn == "" {
		return nil, errors.New("нет DSN")
	}

	gormDb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	sqlDb, err := gormDb.DB()
	if err != nil {
		return nil, fmt.Errorf("не удалось получить объект sql.DB: %v", err)
	}

	sqlDb.SetMaxOpenConns(10)
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetConnMaxLifetime(1 * time.Minute)

	if err = sqlDb.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping failed: %w", err)
	}
	return &Db{gormDb}, nil
}

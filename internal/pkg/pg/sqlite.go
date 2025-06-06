package pg

import (
	"errors"
	"fmt"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
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
		return nil, fmt.Errorf("error connecting to a database: %w", err)
	}

	sqlDb, err := gormDb.DB()
	if err != nil {
		return nil, fmt.Errorf("failing to get sql.DB: %w", err)
	}

	sqlDb.SetMaxOpenConns(10)
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetConnMaxLifetime(1 * time.Minute)

	if err = sqlDb.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping failed: %w", err)
	}
	return &Db{gormDb}, nil
}

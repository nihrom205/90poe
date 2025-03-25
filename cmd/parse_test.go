package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/nihrom205/90poe/internal/app/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func initDb(t *testing.T) *gorm.DB {
	dsn := "file::memory:?cache=shared"
	path := "file://.././internal/app/migrations"

	sqlDB, err := sql.Open("sqlite", dsn)
	require.NoError(t, err, "Failed to open SQLite DB")
	defer sqlDB.Close()

	// Инициализация драйвера для SQLite
	driver, err := sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
	require.NoError(t, err, "Failed to init GORM")

	m, err := migrate.NewWithDatabaseInstance(
		path,
		"sqlite3",
		driver)

	require.NoError(t, err, "Failed to migration not created")

	err = m.Up()
	if err != nil && errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("Migration not executed")
	}

	gormDb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err, "Failed to open GORM DB")

	return gormDb
}

func initData(db *gorm.DB) {
	db.Create(&repository.Port{
		Model:       gorm.Model{},
		Key:         "AEAUH",
		Name:        "Abu Dhabi",
		City:        "Abu Dhabi",
		Country:     "United Arab Emirates",
		Alias:       []string{},
		Regions:     []string{},
		Coordinates: []float64{54.37, 24.47},
		Province:    "Abu Z¸aby [Abu Dhabi]",
		Timezone:    "Asia/Dubai",
		Unlocs:      []string{"AEAUH"},
		Code:        "52001",
	})
}

func TestGetPortByKeySuccess(t *testing.T) {
	// Устанавливаем переменные окружения
	_ = os.Setenv("DSN", "file::memory:?cache=shared")
	_ = os.Setenv("HTTP_ADDR", ":8080")
	_ = os.Setenv("MIGRATIONS_PATH", "file://.././internal/app/migrations")

	db := initDb(t)
	initData(db)

	go main()

	// Отправляем GET-запрос
	resp, err := http.Get("http://localhost:8080/port/AEAUH")

	assert.NoError(t, err)
	defer resp.Body.Close()

	// Проверяем статус-код
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	//Проверяем содержимое ответа
	var port *repository.Port
	err = json.NewDecoder(resp.Body).Decode(&port)
	assert.NoError(t, err)
	assert.Equal(t, port.ID, uint(0))
	assert.Equal(t, port.Key, "AEAUH")
	assert.Equal(t, port.Name, "Abu Dhabi")
}

func TestGetPortByKeyNotFound(t *testing.T) {
	// Устанавливаем переменные окружения
	_ = os.Setenv("DSN", "file::memory:?cache=shared")
	_ = os.Setenv("HTTP_ADDR", ":8080")
	_ = os.Setenv("MIGRATIONS_PATH", "file://.././internal/app/migrations")

	db := initDb(t)
	initData(db)

	go main()

	// Отправляем GET-запрос
	resp, err := http.Get("http://localhost:8080/port/fail")

	assert.NoError(t, err)
	defer resp.Body.Close()

	// Проверяем статус-код
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

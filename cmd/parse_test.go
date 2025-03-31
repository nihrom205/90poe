package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nihrom205/90poe/internal/pkg/pg"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/nihrom205/90poe/internal/app/service"
	"github.com/nihrom205/90poe/internal/app/transport/httpserver"
	"github.com/stretchr/testify/require"

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

func initHttpServer(db *gorm.DB) *httpserver.HttpServer {
	logger := getLogger()
	portRepo := repository.NewPortRepository(&pg.Db{DB: db})
	portService := service.NewPortService(portRepo, &logger)
	return httpserver.NewHttpServer(portService)
}

func TestGetPortSuccess(t *testing.T) {
	db := initDb(t)
	initData(db)
	server := initHttpServer(db)
	router := mux.NewRouter()

	router.HandleFunc("/port/{key}", server.GetPort).Methods(http.MethodGet)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Отправляем GET-запрос
	var netClient = &http.Client{Timeout: time.Second * 10}
	resp, err := netClient.Get(ts.URL + "/port/AEAUH")

	assert.NoError(t, err)
	defer resp.Body.Close()

	// Проверяем статус-код
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	//Проверяем содержимое ответа
	var port repository.Port
	err = json.NewDecoder(resp.Body).Decode(&port)
	assert.NoError(t, err)
	assert.Equal(t, port.ID, uint(0))
	assert.Equal(t, port.Key, "AEAUH")
	assert.Equal(t, port.Name, "Abu Dhabi")
}

func TestGetPortNotFound(t *testing.T) {
	db := initDb(t)
	initData(db)
	server := initHttpServer(db)
	router := mux.NewRouter()

	router.HandleFunc("/port/{key}", server.GetPort).Methods(http.MethodGet)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Отправляем GET-запрос
	var netClient = &http.Client{Timeout: time.Second * 10}
	resp, err := netClient.Get(ts.URL + "/port/fail")

	assert.NoError(t, err)
	defer resp.Body.Close()

	// Проверяем статус-код
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

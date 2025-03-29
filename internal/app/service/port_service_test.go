package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/glebarez/sqlite"
	"github.com/nihrom205/90poe/internal/app/repository"
	"github.com/nihrom205/90poe/internal/pkg"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func initTest() (*gorm.DB, sqlmock.Sqlmock, error) {
	database, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, errors.New("не удалось подключить mock db")
	}

	// Ожидаем запрос `SELECT sqlite_version()`
	rows := sqlmock.NewRows([]string{"version"}).
		AddRow("3.36.0")
	mock.ExpectQuery("select sqlite_version()").
		WillReturnRows(rows)

	gormDb, err := gorm.Open(sqlite.Dialector{Conn: database}, &gorm.Config{})
	if err != nil {
		return nil, nil, errors.New("не удалось подключить GORM")
	}
	return gormDb, mock, nil
}

func TestPortServiceGetPortByKeySuccess(t *testing.T) {
	logger := getLogger()
	gormDb, mock, err := initTest()
	if err != nil {
		t.Fatal(err)
		return
	}
	db, _ := gormDb.DB()
	defer db.Close()

	// Ожидаемый запрос и результат
	rows := mock.NewRows([]string{"key", "name"}).
		AddRow("AEAUH", "Abu Dhabi")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	// Вызываем тестируемую функцию
	portRepo := repository.NewPortRepository(&pkg.Db{
		DB: gormDb,
	})
	ctx := context.Background()
	service := NewPortService(portRepo, &logger)
	port, err := service.GetPortByKey(ctx, "AEAUH")

	assert.NoError(t, err)
	assert.NotNil(t, port)
	assert.Equal(t, "AEAUH", port.Key)
	assert.Equal(t, "Abu Dhabi", port.Name)

	// Проверяем, что все ожидания выполнены
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("не выполнены ожидания: %v", err)
	}
}

func TestPortServiceGetPortByKeyNotFound(t *testing.T) {
	logger := getLogger()
	gormDb, mock, err := initTest()
	if err != nil {
		t.Fatal(err)
		return
	}
	db, _ := gormDb.DB()
	defer db.Close()

	// Ожидаемый запрос и результат
	rows := mock.NewRows([]string{"key", "name"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	// Вызываем тестируемую функцию
	portRepo := repository.NewPortRepository(&pkg.Db{
		DB: gormDb,
	})
	ctx := context.Background()
	service := NewPortService(portRepo, &logger)
	port, err := service.GetPortByKey(ctx, "FakeKey")

	assert.Error(t, err)
	assert.Nil(t, port)

	// Проверяем, что все ожидания выполнены
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("не выполнены ожидания: %v", err)
	}
}

func TestPortServiceParsing(t *testing.T) {
	logger := getLogger()
	// Создаём mock-репозиторий
	mockRepo := new(MockPortRepository)

	// Настраиваем ожидание для GetPortByKey
	port := getPort()

	mockRepo.On("GetPortByKey", port.Key).
		Return((*repository.Port)(nil), gorm.ErrRecordNotFound)

	// Настраиваем ожидание для CreatePort
	mockRepo.On("CreatePort", port).
		Return(port, nil)

	// Создаём сервис с mock-репозиторием
	service := NewPortService(mockRepo, &logger)

	// Вызываем метод SavePort
	service.ProcessingJson(context.Background(), io.NopCloser(strings.NewReader(getJsonData())))

	// Проверяем, что методы были вызваны с ожидаемыми аргументами
	mockRepo.AssertCalled(t, "GetPortByKey", port.Key)
	//mockRepo.AssertCalled(t, "CreatePort", &port)

	// Проверяем, что все ожидания выполнены
	mockRepo.AssertExpectations(t)
}

func getJsonData() string {
	return `{
		"AEAJM": {
			"name": "Ajman",
			"city": "Ajman",
			"country": "United Arab Emirates",
			"alias": [],
			"regions": [],
			"coordinates": [
			  55.5136433,
			  25.4052165
			],
			"province": "Ajman",
			"timezone": "Asia/Dubai",
			"unlocs": [
			  "AEAJM"
			],
			"code": "52000"
		  }
		}`
}

func getPort() *repository.Port {
	jsonData2 := `{
			"name": "Ajman",
			"city": "Ajman",
			"country": "United Arab Emirates",
			"alias": [],
			"regions": [],
			"coordinates": [
			  55.5136433,
			  25.4052165
			],
			"province": "Ajman",
			"timezone": "Asia/Dubai",
			"unlocs": [
			  "AEAJM"
			],
			"code": "52000"
		  }`

	location := Location{}
	decoder := json.NewDecoder(strings.NewReader(jsonData2))
	err := decoder.Decode(&location)
	if err != nil {
		log.Fatalf("PortService - getPort: Ошибка при декодировании объекта: %v", err)
	}

	port, _ := mapperToDB(keyAndLocation{
		Key: "AEAJM",
		Loc: location,
	})

	return port
}

func getLogger() zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006-01-02T15:04:05Z07:00"}).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()
}

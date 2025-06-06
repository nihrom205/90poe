package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/nihrom205/90poe/internal/pkg/pg"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nihrom205/90poe/internal/app/config"
	"github.com/nihrom205/90poe/internal/app/repository"
	"github.com/nihrom205/90poe/internal/app/service"
	"github.com/nihrom205/90poe/internal/app/transport/httpserver"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func run() error {
	logger := getLogger()
	cfg := config.Read()

	db, err := pg.NewDb(cfg.DSN)
	if err != nil {
		return fmt.Errorf("pg.Db failed: %w", err)
	}

	// Migration
	if db != nil {
		logger.Info().Msg("Start Sqlite migrations")
		if err = runSqliteMigrations(db, cfg.MigrationsPath); err != nil {
			return fmt.Errorf("runSqliteMigrations failed: %w", err)
		}
	}
	sqlDb, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed connection to a database: %w", err)
	}
	defer func(sqlDb *sql.DB) {
		err := sqlDb.Close()
		if err != nil {
			logger.Error().Err(err).Msg("failed to close database connection")
		}
	}(sqlDb)

	// Create Repositories
	portRepo := repository.NewPortRepository(db)

	// Create Services
	portService := service.NewPortService(portRepo, &logger)

	// Create HttpServer
	httpServer := httpserver.NewHttpServer(portService)

	// Create http router
	router := mux.NewRouter()

	router.HandleFunc("/ports", httpServer.LoadPorts).Methods(http.MethodPost)
	router.HandleFunc("/port/{key}", httpServer.GetPort).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}

	// Listen OS signal
	stopped := make(chan struct{})
	go func() {
		signCh := make(chan os.Signal, 1)
		signal.Notify(signCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-signCh
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err = srv.Shutdown(ctx); err != nil {
			logger.Error().Err(err).Msg("HTTP Server Shutdown Error")
		}
		close(stopped)
	}()

	// Start http server  запущен
	logger.Info().Msgf("Starting Http server on %s port", cfg.HTTPAddr)
	if err = srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP Server Shutdown Error: %w", err)
	}

	<-stopped

	logger.Info().Msg("Http server completed his work!")

	return nil
}

func runSqliteMigrations(db *pg.Db, path string) error {
	if path == "" {
		return errors.New("no migrations path provided")
	}

	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed connection to a database: %w", err)
	}

	// Инициализация драйвера для SQLite
	driver, err := sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
	if err != nil {
		return errors.New("failed initialization of SQLite driver")
	}

	m, err := migrate.NewWithDatabaseInstance(
		path,
		"sqlite3",
		driver)

	if err != nil {
		return fmt.Errorf("instance migration not created: %s", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration not executed: %s", err)
	}
	return nil
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

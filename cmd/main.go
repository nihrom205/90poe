package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nihrom205/90poe/internal/app/config"
	"github.com/nihrom205/90poe/internal/app/repository"
	"github.com/nihrom205/90poe/internal/app/service"
	"github.com/nihrom205/90poe/internal/app/transport/httpserver"
	"github.com/nihrom205/90poe/internal/pkg"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func run() error {
	cfg := config.Read()

	db, err := pkg.NewDb(cfg.DSN)
	if err != nil {
		return fmt.Errorf("pg.Db failed: %w", err)
	}

	// Migration
	if db != nil {
		log.Println("Start Sqlite migrations")
		if err = runSqliteMigrations(cfg.DSN, cfg.MigrationsPath); err != nil {
			return fmt.Errorf("runSqliteMigrations failed: %w", err)
		}
	}

	// Create Repositories
	portRepo := repository.NewPortRepository(db)

	// Create Services
	portService := service.NewPortService(portRepo)

	// Create HttpServer
	httpServer := httpserver.NewHttpServer(portService)

	// Create http router
	router := mux.NewRouter()

	router.HandleFunc("/ports", httpServer.Processing).Methods(http.MethodPost)
	router.HandleFunc("/port/{key}", httpServer.GetPortByKey).Methods(http.MethodGet)

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
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(stopped)
	}()

	// Start http server  запущен
	log.Printf("Starting Http server on %s port", cfg.HTTPAddr)
	if err = srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe stoped with error: %v", err)
	}

	<-stopped

	log.Printf("Http server completed his work!")

	return nil
}

func runSqliteMigrations(dsn, path string) error {
	if path == "" {
		return errors.New("no migrations path provided")
	}
	if dsn == "" {
		return errors.New("no DSN provided")
	}

	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return fmt.Errorf("failed connection to a database: %v", err)
	}
	defer sqlDB.Close()

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

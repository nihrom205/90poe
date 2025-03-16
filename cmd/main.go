package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nihrom205/90poe/internal/app/config"
	"github.com/nihrom205/90poe/internal/pkg"
	"github.com/nihrom205/90poe/internal/transport/httpserver"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func run() error {
	cfg := config.Read()

	_, err := pkg.NewDb(&cfg)
	if err != nil {
		return fmt.Errorf("pg.Db failed: %w", err)
	}

	// Create Repositories

	// Create Services

	// Create HttpServer
	httpServer := httpserver.NewHttpServer()

	// Create http router
	router := mux.NewRouter()

	router.HandleFunc("/test", httpServer.MyTestHandler).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}

	// listen OS signal
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

	// start http server
	log.Printf("Starting HTTP server on %s", cfg.HTTPAddr)
	if err = srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-stopped

	log.Printf("The server completed the work!")

	return nil
}

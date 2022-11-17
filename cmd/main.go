package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/KarolosLykos/json-validation-service/internal/api"
	"github.com/KarolosLykos/json-validation-service/internal/api/server"
	"github.com/KarolosLykos/json-validation-service/internal/config"
	"github.com/KarolosLykos/json-validation-service/internal/logger/logruslog"
	"github.com/KarolosLykos/json-validation-service/internal/storage"
	"github.com/KarolosLykos/json-validation-service/internal/storage/postgres"
)

func main() {
	if err := start(); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
}

func start() error {
	ctx := context.TODO()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log := logruslog.DefaultLogger(cfg.Debug)

	p := postgres.New(cfg, log)

	db, err := p.Connect(ctx)
	if err != nil {
		return err
	}

	if err := db.Initialize(ctx); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	s := server.New(ctx, cfg, log, db)
	s.Start(ctx)

	event := <-quit

	log.Debug(ctx, fmt.Sprintf("received signal: %v", event))

	return shutdown(ctx, s, db)
}

func shutdown(ctx context.Context, s api.API, db storage.Storage) error {
	s.Shutdown(ctx)

	return db.Shutdown(ctx)
}

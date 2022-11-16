package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/KarolosLykos/json-validation-service/internal/api/server"
	"github.com/KarolosLykos/json-validation-service/internal/config"
	"github.com/KarolosLykos/json-validation-service/internal/logger/logruslog"
)

func main() {
	ctx := context.TODO()

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	
	log := logruslog.DefaultLogger(cfg.Debug)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	s := server.New(ctx, cfg, log)
	s.Start(ctx)

	event := <-quit

	log.Debug(ctx, fmt.Sprintf("received signal: %v", event))

	s.Shutdown(ctx)
}

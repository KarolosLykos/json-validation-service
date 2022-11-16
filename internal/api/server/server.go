package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"

	"github.com/KarolosLykos/json-validation-service/internal/api"
	"github.com/KarolosLykos/json-validation-service/internal/api/server/routes"
	"github.com/KarolosLykos/json-validation-service/internal/config"
	"github.com/KarolosLykos/json-validation-service/internal/logger"
)

type server struct {
	cfg     *config.Config
	log     logger.Logger
	server  *http.Server
	handler http.Handler
}

func New(ctx context.Context, cfg *config.Config, log logger.Logger) api.API {
	log.Debug(ctx, "create new server")

	corsOptions := []handlers.CORSOption{
		handlers.AllowedMethods([]string{http.MethodPost, http.MethodGet}),
		handlers.AllowedHeaders([]string{"content-type"}),
	}

	router := routes.SetupRoutes(ctx, log)

	handler := handlers.CORS(corsOptions...)(router)

	return &server{
		cfg:     cfg,
		log:     log,
		handler: handler,
	}
}

func (s *server) Start(ctx context.Context) {
	s.log.Debug(ctx, "starting server")

	addr := fmt.Sprintf("%s:%s", s.cfg.HTTP.IP, s.cfg.HTTP.Port)

	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.handler,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 15,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				s.log.Panic(ctx, err, "could not initiate https server", s.cfg.HTTP.Port)
			}
		}
	}()

	s.log.Info(ctx, "server started on ", addr)
}

func (s *server) Shutdown(ctx context.Context) {
	s.log.Debug(ctx, "shutting down server")

	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Error(ctx, err, "could not shutdown server")
	}
}

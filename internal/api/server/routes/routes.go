package routes

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/KarolosLykos/json-validation-service/internal/api/server/handlers"
	"github.com/KarolosLykos/json-validation-service/internal/api/server/middleware"
	"github.com/KarolosLykos/json-validation-service/internal/logger"
)

func SetupRoutes(ctx context.Context, log logger.Logger) http.Handler {
	log.Debug(ctx, "setting up routes")

	router := mux.NewRouter().StrictSlash(true)

	m := middleware.New(log)

	router.Use(m.RecoverPanic)

	h := handlers.New(log)

	router.HandleFunc("/schema/{schemaID}", h.Upload()).Methods(http.MethodPost)
	router.HandleFunc("/schema/{schemaID}", h.Download()).Methods(http.MethodGet)
	router.HandleFunc("/validate/{schemaID}", h.Validate()).Methods(http.MethodPost)

	return router
}

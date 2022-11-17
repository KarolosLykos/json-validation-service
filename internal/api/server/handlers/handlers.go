package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/KarolosLykos/json-validation-service/internal/logger"
	"github.com/KarolosLykos/json-validation-service/internal/storage"
)

type Handler struct {
	log logger.Logger
	db  storage.Storage
}

func New(log logger.Logger, db storage.Storage) *Handler {
	return &Handler{
		log: log,
		db:  db,
	}
}

func (h *Handler) Upload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		schemaID := vars["schemaID"]

		fmt.Println(schemaID)
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			return
		}

		fmt.Println(string(body))

		h.log.Debug(ctx, "upload")
	}
}

func (h *Handler) Download() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		schemaID := vars["schemaID"]

		fmt.Println(schemaID)

		h.log.Debug(ctx, "download")
	}
}

func (h *Handler) Validate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		schemaID := vars["schemaID"]

		fmt.Println(schemaID)

		h.log.Debug(ctx, "validate")
	}
}

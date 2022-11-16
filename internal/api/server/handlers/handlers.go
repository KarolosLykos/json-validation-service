package handlers

import (
	"net/http"

	"github.com/KarolosLykos/json-validation-service/internal/logger"
)

type Handler struct {
	log logger.Logger
}

func New(log logger.Logger) *Handler {
	return &Handler{
		log: log,
	}
}

func (h *Handler) Upload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.log.Debug(r.Context(), "upload")
	}
}

func (h *Handler) Download() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.log.Debug(r.Context(), "download")
	}
}

func (h *Handler) Validate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.log.Debug(r.Context(), "validate")
	}
}

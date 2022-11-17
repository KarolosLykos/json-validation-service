package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/KarolosLykos/json-validation-service/internal/logger"
	"github.com/KarolosLykos/json-validation-service/internal/service"
	"github.com/KarolosLykos/json-validation-service/internal/utils/exceptions"
)

type Handler struct {
	log logger.Logger
	srv service.Service
}

func New(log logger.Logger, srv service.Service) *Handler {
	return &Handler{
		log: log,
		srv: srv,
	}
}

type Response struct {
	Action  string      `json:"action"`
	ID      string      `json:"id"`
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

func (h *Handler) Upload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		schemaID := vars["schemaID"]

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			responseError(w, "uploadSchema", schemaID, err)

			return
		}

		if err = h.srv.UploadSchema(ctx, schemaID, string(body)); err != nil {
			responseError(w, "uploadSchema", schemaID, err)

			return
		}

		responseSuccess(w, http.StatusCreated, "uploadSchema", schemaID, nil)
	}
}

func (h *Handler) Download() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		schemaID := vars["schemaID"]

		fmt.Println(schemaID)

		s, err := h.srv.DownloadSchema(ctx, schemaID)
		if err != nil {
			responseError(w, "downloadSchema", schemaID, err)

			return
		}

		responseSuccess(w, http.StatusOK, "downloadSchema", schemaID, s)
	}
}

func (h *Handler) Validate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		schemaID := vars["schemaID"]

		payload := make(map[string]interface{})

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			responseError(w, "validateSchema", schemaID, err)

			return
		}

		if err := h.srv.ValidateSchema(ctx, schemaID, payload); err != nil {
			responseError(w, "validateSchema", schemaID, err)

			return
		}

		responseSuccess(w, http.StatusOK, "validateSchema", schemaID, nil)
	}
}

func responseError(w http.ResponseWriter, action, schemaID string, errMsg error) {
	statusCode := http.StatusInternalServerError

	if errors.Is(errMsg, exceptions.ErrInvalidJSON) || errors.Is(errMsg, exceptions.ErrNotFound) {
		statusCode = http.StatusBadRequest
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	res := &Response{
		Action:  action,
		ID:      schemaID,
		Status:  "error",
		Message: errMsg.Error(),
	}

	payload, err := json.Marshal(res)
	if err != nil {
		http.Error(w, exceptions.ErrInternalServerError.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(payload)
	if err != nil {
		http.Error(w, exceptions.ErrInternalServerError.Error(), http.StatusInternalServerError)
	}
}

func responseSuccess(w http.ResponseWriter, statusCode int, action, schemaID string, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	res := &Response{
		Action: action,
		ID:     schemaID,
		Status: "success",
	}

	if payload != nil {
		res.Payload = payload
	}

	p, err := json.Marshal(res)
	if err != nil {
		http.Error(w, exceptions.ErrInternalServerError.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(p)
	if err != nil {
		http.Error(w, exceptions.ErrInternalServerError.Error(), http.StatusInternalServerError)
	}
}

package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KarolosLykos/json-validation-service/internal/api/server/handlers"
	"github.com/KarolosLykos/json-validation-service/internal/config"
	"github.com/KarolosLykos/json-validation-service/internal/logger/logruslog"
	"github.com/KarolosLykos/json-validation-service/internal/service"
	mock_service "github.com/KarolosLykos/json-validation-service/internal/service/mock"
	"github.com/KarolosLykos/json-validation-service/internal/utils/exceptions"
)

func TestHandler_Download(t *testing.T) {
	tc := []struct {
		name        string
		schemaID    string
		serviceStub func(srv *mock_service.MockService)
		payload     string
		statusCode  int
		res         *handlers.Response
	}{
		{
			name:     "status created",
			schemaID: "config-schema",
			payload:  "",
			serviceStub: func(srv *mock_service.MockService) {
				srv.EXPECT().
					UploadSchema(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			statusCode: http.StatusCreated,
			res: &handlers.Response{
				Action: "uploadSchema",
				ID:     "config-schema",
				Status: "success",
			},
		},
		{
			name:     "bad request",
			schemaID: "config-schema",
			payload:  "",
			serviceStub: func(srv *mock_service.MockService) {
				srv.EXPECT().
					UploadSchema(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(exceptions.ErrAlreadyExists)
			},
			statusCode: http.StatusConflict,
			res: &handlers.Response{
				Action:  "uploadSchema",
				ID:      "config-schema",
				Status:  "error",
				Message: exceptions.ErrAlreadyExists.Error(),
			},
		},
		{
			name:     "internal server error",
			schemaID: "config-schema",
			payload:  "",
			serviceStub: func(srv *mock_service.MockService) {
				srv.EXPECT().
					UploadSchema(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(exceptions.ErrCreateSchema)
			},
			statusCode: http.StatusInternalServerError,
			res: &handlers.Response{
				Action:  "uploadSchema",
				ID:      "config-schema",
				Status:  "error",
				Message: exceptions.ErrCreateSchema.Error(),
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srv := mock_service.NewMockService(ctrl)
			tt.serviceStub(srv)

			h := helperNewHandler(t, srv)

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/schema/schemaID:%s", tt.schemaID), bytes.NewBuffer([]byte(tt.payload)))
			r = mux.SetURLVars(r, map[string]string{"schemaID": tt.schemaID})

			h.Upload()(w, r)

			require.Equal(t, tt.statusCode, w.Code)

			res := &handlers.Response{}
			err := json.NewDecoder(w.Body).Decode(res)
			require.NoError(t, err)

			assert.Equal(t, tt.res, res)
		})
	}
}

func TestHandler_Upload(t *testing.T) {
	tc := []struct {
		name        string
		schemaID    string
		serviceStub func(srv *mock_service.MockService)
		schema      string
		statusCode  int
		res         *handlers.Response
	}{
		{
			name:     "status ok",
			schemaID: "config-schema",
			schema:   `{"valid":"schema"}`,
			serviceStub: func(srv *mock_service.MockService) {
				srv.EXPECT().
					DownloadSchema(gomock.Any(), gomock.Any()).
					Times(1).
					Return(`{"valid":"schema"}`, nil)
			},
			statusCode: http.StatusOK,
			res: &handlers.Response{
				Action:  "downloadSchema",
				ID:      "config-schema",
				Status:  "success",
				Payload: `{"valid":"schema"}`,
			},
		},
		{
			name:     "not found",
			schemaID: "config-schema",
			schema:   "",
			serviceStub: func(srv *mock_service.MockService) {
				srv.EXPECT().
					DownloadSchema(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", exceptions.ErrNotFound)
			},
			statusCode: http.StatusBadRequest,
			res: &handlers.Response{
				Action:  "downloadSchema",
				ID:      "config-schema",
				Status:  "error",
				Message: exceptions.ErrNotFound.Error(),
			},
		},
		{
			name:     "internal server error",
			schemaID: "config-schema",
			schema:   "",
			serviceStub: func(srv *mock_service.MockService) {
				srv.EXPECT().
					DownloadSchema(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", exceptions.ErrDownloadSchema)
			},
			statusCode: http.StatusInternalServerError,
			res: &handlers.Response{
				Action:  "downloadSchema",
				ID:      "config-schema",
				Status:  "error",
				Message: exceptions.ErrDownloadSchema.Error(),
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srv := mock_service.NewMockService(ctrl)
			tt.serviceStub(srv)

			h := helperNewHandler(t, srv)

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/schema/schemaID:%s", tt.schemaID), nil)
			r = mux.SetURLVars(r, map[string]string{"schemaID": tt.schemaID})

			h.Download()(w, r)

			require.Equal(t, tt.statusCode, w.Code)

			res := &handlers.Response{}
			err := json.NewDecoder(w.Body).Decode(res)
			require.NoError(t, err)

			assert.Equal(t, tt.res, res)
		})
	}
}

func TestHandler_Validate(t *testing.T) {
	tc := []struct {
		name        string
		schemaID    string
		serviceStub func(srv *mock_service.MockService)
		payload     string
		statusCode  int
		res         *handlers.Response
	}{
		{
			name:     "status ok",
			schemaID: "config-schema",
			payload:  `{"valid":"schema"}`,
			serviceStub: func(srv *mock_service.MockService) {
				srv.EXPECT().
					ValidateSchema(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			statusCode: http.StatusOK,
			res: &handlers.Response{
				Action: "validateSchema",
				ID:     "config-schema",
				Status: "success",
			},
		},
		{
			name:     "invalid payload",
			schemaID: "config-schema",
			payload:  "",
			serviceStub: func(srv *mock_service.MockService) {
				srv.EXPECT().
					ValidateSchema(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0).
					Return(nil)
			},
			statusCode: http.StatusBadRequest,
			res: &handlers.Response{
				Action:  "validateSchema",
				ID:      "config-schema",
				Status:  "error",
				Message: io.EOF.Error(),
			},
		},
		{
			name:     "not found",
			schemaID: "config-schema",
			payload:  "{}",
			serviceStub: func(srv *mock_service.MockService) {
				srv.EXPECT().
					ValidateSchema(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(exceptions.ErrNotFound)
			},
			statusCode: http.StatusBadRequest,
			res: &handlers.Response{
				Action:  "validateSchema",
				ID:      "config-schema",
				Status:  "error",
				Message: exceptions.ErrNotFound.Error(),
			},
		},
		{
			name:     "not valid payload",
			schemaID: "config-schema",
			payload:  "{}",
			serviceStub: func(srv *mock_service.MockService) {
				srv.EXPECT().
					ValidateSchema(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(exceptions.ErrValidation)
			},
			statusCode: http.StatusBadRequest,
			res: &handlers.Response{
				Action:  "validateSchema",
				ID:      "config-schema",
				Status:  "error",
				Message: exceptions.ErrValidation.Error(),
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srv := mock_service.NewMockService(ctrl)
			tt.serviceStub(srv)

			h := helperNewHandler(t, srv)

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/schema/schemaID:%s", tt.schemaID), bytes.NewBuffer([]byte(tt.payload)))
			r = mux.SetURLVars(r, map[string]string{"schemaID": tt.schemaID})

			h.Validate()(w, r)

			require.Equal(t, tt.statusCode, w.Code)

			res := &handlers.Response{}
			err := json.NewDecoder(w.Body).Decode(res)
			require.NoError(t, err)

			assert.Equal(t, tt.res, res)
		})
	}
}

func helperNewHandler(t *testing.T, srv service.Service) *handlers.Handler {
	t.Helper()
	cfg, _ := config.Load()

	log := logruslog.DefaultLogger(cfg)

	return handlers.New(log, srv)
}

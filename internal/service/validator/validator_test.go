package validator_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/KarolosLykos/json-validation-service/internal/config"
	"github.com/KarolosLykos/json-validation-service/internal/logger/logruslog"
	"github.com/KarolosLykos/json-validation-service/internal/service"
	"github.com/KarolosLykos/json-validation-service/internal/service/validator"
	"github.com/KarolosLykos/json-validation-service/internal/storage"
	mock_storage "github.com/KarolosLykos/json-validation-service/internal/storage/mock"
	"github.com/KarolosLykos/json-validation-service/internal/utils/exceptions"
)

func TestValidator_UploadSchema(t *testing.T) {
	tc := []struct {
		name      string
		schemaID  string
		schema    string
		storeStub func(store *mock_storage.MockStorage)
		err       error
	}{
		{
			name:     "success",
			schemaID: "config-schema",
			schema:   `{ "valid": "json" }`,
			storeStub: func(store *mock_storage.MockStorage) {
				store.EXPECT().
					CreateSchema(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			err: nil,
		},
		{
			name:     "invalid json",
			schemaID: "config-schema",
			schema:   `{ "invalid"  }`,
			storeStub: func(store *mock_storage.MockStorage) {
				store.EXPECT().
					CreateSchema(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0).
					Return(nil)
			},
			err: exceptions.ErrInvalidJSON,
		},
		{
			name:     "already exists",
			schemaID: "config-schema",
			schema:   `{ "valid": "json" }`,
			storeStub: func(store *mock_storage.MockStorage) {
				store.EXPECT().
					CreateSchema(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(&pq.Error{Code: "23505"})
			},
			err: exceptions.ErrAlreadyExists,
		},
		{
			name:     "generic error",
			schemaID: "config-schema",
			schema:   `{ "valid": "json" }`,
			storeStub: func(store *mock_storage.MockStorage) {
				store.EXPECT().
					CreateSchema(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("error"))
			},
			err: exceptions.ErrCreateSchema,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_storage.NewMockStorage(ctrl)
			tt.storeStub(store)

			v := helperNewValidator(t, store)

			err := v.UploadSchema(ctx, tt.schemaID, tt.schema)
			if tt.err != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidator_DownloadSchema(t *testing.T) {
	tc := []struct {
		name      string
		schemaID  string
		storeStub func(store *mock_storage.MockStorage)
		err       error
	}{
		{
			name:     "success",
			schemaID: "config-schema",
			storeStub: func(store *mock_storage.MockStorage) {
				store.EXPECT().
					GetSchema(gomock.Any(), gomock.Any()).
					Times(1).
					Return(`{ "valid": "json" }`, nil)
			},
			err: nil,
		},
		{
			name:     "not found",
			schemaID: "config-schema",
			storeStub: func(store *mock_storage.MockStorage) {
				store.EXPECT().
					GetSchema(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", gorm.ErrRecordNotFound)
			},
			err: exceptions.ErrNotFound,
		},
		{
			name:     "generic error",
			schemaID: "config-schema",
			storeStub: func(store *mock_storage.MockStorage) {
				store.EXPECT().
					GetSchema(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", errors.New("error"))
			},
			err: exceptions.ErrDownloadSchema,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_storage.NewMockStorage(ctrl)
			tt.storeStub(store)

			v := helperNewValidator(t, store)

			s, err := v.DownloadSchema(ctx, tt.schemaID)
			if tt.err != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
				require.NotNil(t, s)
			}
		})
	}
}

func TestValidator_ValidateSchema(t *testing.T) {
	tc := []struct {
		name      string
		schemaID  string
		payload   map[string]interface{}
		storeStub func(store *mock_storage.MockStorage)
		err       error
	}{
		{
			name:     "success",
			schemaID: "config-schema",
			payload:  map[string]interface{}{"valid": "json"},
			storeStub: func(store *mock_storage.MockStorage) {
				store.EXPECT().
					GetSchema(gomock.Any(), gomock.Any()).
					Times(1).
					Return(`{ "valid": "json" }`, nil)
			},
			err: nil,
		},
		{
			name:     "not found",
			schemaID: "config-schema",
			payload:  map[string]interface{}{"valid": "json"},
			storeStub: func(store *mock_storage.MockStorage) {
				store.EXPECT().
					GetSchema(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", exceptions.ErrNotFound)
			},
			err: exceptions.ErrNotFound,
		},
		{
			name:     "cannot add resource",
			schemaID: "config-schema",
			payload:  map[string]interface{}{"valid": "json"},
			storeStub: func(store *mock_storage.MockStorage) {
				store.EXPECT().
					GetSchema(gomock.Any(), gomock.Any()).
					Times(1).
					Return("-", nil)
			},
			err: exceptions.ErrValidateSchema,
		},
		{
			name:     "cannot validate schema",
			schemaID: "config-schema",
			payload:  map[string]interface{}{"source": "something"},
			storeStub: func(store *mock_storage.MockStorage) {
				store.EXPECT().
					GetSchema(gomock.Any(), gomock.Any()).
					Times(1).
					Return(`{
					  "$schema":"http://json-schema.org/draft-04/schema#",
					  "type": "object",
					  "properties": {
						"source":{
						  "type": "string"
						},
						"destination": {
						  "type": "string"
						}
					  },
					  "required": ["source", "destination"]
					}`, nil)
			},
			err: exceptions.ErrValidation,
		},
		{
			name:     "should remove null values and throw error on required fields",
			schemaID: "config-schema",
			payload:  map[string]interface{}{"source": "value", "destination": nil},
			storeStub: func(store *mock_storage.MockStorage) {
				store.EXPECT().
					GetSchema(gomock.Any(), gomock.Any()).
					Times(1).
					Return(`{
					  "$schema":"http://json-schema.org/draft-04/schema#",
					  "type": "object",
					  "properties": {
						"source":{
						  "type": "string"
						},
						"destination": {
						  "type": "string"
						}
					  },
					  "required": ["source", "destination"]
					}`, nil)
			},
			err: exceptions.ErrValidation,
		},
		{
			name:     "should remove null values on nested objects and throw error on required fields",
			schemaID: "config-schema",
			payload:  map[string]interface{}{"chunks": map[string]interface{}{"number": nil}},
			storeStub: func(store *mock_storage.MockStorage) {
				store.EXPECT().
					GetSchema(gomock.Any(), gomock.Any()).
					Times(1).
					Return(`{
					  "$schema":"http://json-schema.org/draft-04/schema#",
					  "type": "object",
					  "properties": {
						"chunks": {
						  "type": "object",
						  "properties": {
							"size": {
							  "type": "integer"
							},
							"number": {
							  "type": "integer"
							}
						  },
						  "required": ["size"]
						}
					  }
					}`, nil)
			},
			err: exceptions.ErrValidation,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_storage.NewMockStorage(ctrl)
			tt.storeStub(store)

			v := helperNewValidator(t, store)

			err := v.ValidateSchema(ctx, tt.schemaID, tt.payload)
			if tt.err != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func helperNewValidator(t *testing.T, store storage.Storage) service.Service {
	t.Helper()

	cfg, _ := config.Load()
	log := logruslog.DefaultLogger(cfg)

	return validator.New(log, store)
}

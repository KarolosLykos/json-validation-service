package validator_test

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/KarolosLykos/json-validation-service/internal/config"
	"github.com/KarolosLykos/json-validation-service/internal/logger/logruslog"
	"github.com/KarolosLykos/json-validation-service/internal/service"
	"github.com/KarolosLykos/json-validation-service/internal/service/validator"
	mock_storage "github.com/KarolosLykos/json-validation-service/internal/storage/mock"
)

func TestValidator_UploadSchema(t *testing.T) {
	// db, v := helperNewValidator(t)
}

func TestValidator_DownloadSchema(t *testing.T) {
	// db, v := helperNewValidator(t)
}

func TestValidator_ValidateSchema(t *testing.T) {
	// db, v := helperNewValidator(t)

}

func helperNewValidator(t *testing.T) (*mock_storage.MockStorage, service.Service) {
	t.Helper()

	cfg, _ := config.Load()

	ctrl := gomock.NewController(t)
	db := mock_storage.NewMockStorage(ctrl)

	log := logruslog.DefaultLogger(cfg)

	return db, validator.New(log, db)
}

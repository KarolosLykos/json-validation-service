package validator

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/lib/pq"
	"github.com/santhosh-tekuri/jsonschema"
	"gorm.io/gorm"

	"github.com/KarolosLykos/json-validation-service/internal/logger"
	"github.com/KarolosLykos/json-validation-service/internal/service"
	"github.com/KarolosLykos/json-validation-service/internal/storage"
	"github.com/KarolosLykos/json-validation-service/internal/utils/exceptions"
)

type Validator struct {
	log logger.Logger
	db  storage.Storage
}

func New(log logger.Logger, db storage.Storage) service.Service {
	return &Validator{
		log: log,
		db:  db,
	}
}

func (v *Validator) UploadSchema(ctx context.Context, schemaID, schema string) error {
	v.log.Debug(ctx, "Validator: uploading schema")

	var empty struct{}
	if err := json.Unmarshal([]byte(schema), &empty); err != nil {
		return exceptions.ErrInvalidJSON
	}

	if err := v.db.CreateSchema(ctx, schemaID, schema); err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			return fmt.Errorf("%w:%v", exceptions.ErrAlreadyExists, err)
		}

		return fmt.Errorf("%w:%v", exceptions.ErrCreateSchema, err)
	}

	return nil
}

func (v *Validator) DownloadSchema(ctx context.Context, schemaID string) (string, error) {
	v.log.Debug(ctx, "Validator: downloading schema")

	s, err := v.db.GetSchema(ctx, schemaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", exceptions.ErrNotFound
		}

		return "", fmt.Errorf("%w:%v", exceptions.ErrDownloadSchema, err)
	}

	return s, nil
}

func (v *Validator) ValidateSchema(ctx context.Context, schemaID string, payload map[string]interface{}) error {
	v.log.Debug(ctx, "Validator: validating schema")

	removeNulls(payload)

	payloadB, _ := json.Marshal(payload)

	s, err := v.db.GetSchema(ctx, schemaID)
	if err != nil {
		return err
	}

	compiler := jsonschema.NewCompiler()

	if err = compiler.AddResource(schemaID, strings.NewReader(s)); err != nil {
		return fmt.Errorf("%w:%v", exceptions.ErrValidateSchema, err)
	}

	schema, err := compiler.Compile(schemaID)
	if err != nil {
		return err
	}

	if err = schema.Validate(bytes.NewReader(payloadB)); err != nil {
		return formatValidationError(err)
	}

	return nil
}

// removeNulls https://gist.github.com/ribice/074ad38d9f2fc5c88b20663659988d19.
func removeNulls(m map[string]interface{}) {
	val := reflect.ValueOf(m)
	for _, e := range val.MapKeys() {
		v := val.MapIndex(e)
		if v.IsNil() {
			delete(m, e.String())
			continue
		}
		//nolint:gocritic // type switch
		switch t := v.Interface().(type) {
		// If key is a JSON object (Go Map), use recursion to go deeper
		case map[string]interface{}:
			removeNulls(t)
		}
	}
}

func formatValidationError(err error) error {
	if ve, ok := err.(*jsonschema.ValidationError); ok {
		msg := ve.Message
		for _, c := range ve.Causes {
			msg += c.Message + ","
		}

		return fmt.Errorf("%v:%v", exceptions.ErrValidation, msg)
	}

	return err
}

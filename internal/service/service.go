package service

import (
	"context"
)

type Service interface {
	UploadSchema(ctx context.Context, schemaID, schema string) error
	DownloadSchema(ctx context.Context, schemaID string) (string, error)
	ValidateSchema(ctx context.Context, schemaID string, payload map[string]interface{}) error
}

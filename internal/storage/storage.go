package storage

import (
	"context"
)

type Storage interface {
	Connect(ctx context.Context) (Storage, error)
	Shutdown(ctx context.Context) error
	Initialize(ctx context.Context) error

	CreateSchema(ctx context.Context, schemaID, schema string) error
	GetSchema(ctx context.Context, schemaID string) (string, error)
}

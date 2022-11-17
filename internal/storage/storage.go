package storage

import (
	"context"
)

type Storage interface {
	Connect(ctx context.Context) (Storage, error)
	Shutdown(ctx context.Context) error
	Initialize(ctx context.Context) error
}

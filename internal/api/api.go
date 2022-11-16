package api

import (
	"context"
)

type API interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context)
}

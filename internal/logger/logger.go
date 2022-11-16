package logger

import (
	"context"
)

type Logger interface {
	Panic(ctx context.Context, err error, messages ...interface{})
	Error(ctx context.Context, err error, messages ...interface{})
	Info(ctx context.Context, messages ...interface{})
	Debug(ctx context.Context, messages ...interface{})
}

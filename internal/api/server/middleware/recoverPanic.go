package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/KarolosLykos/json-validation-service/internal/logger"
)

var ErrRecover = errors.New("recovering from error")

type Middleware struct {
	log logger.Logger
}

func New(log logger.Logger) *Middleware {
	return &Middleware{
		log: log,
	}
}

// RecoverPanic handles any panic that may occur.
func (m *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		defer func() {
			if err := recover(); err != nil {
				m.log.Error(ctx, fmt.Errorf("%w: %v", ErrRecover, err), "middleware recovering from panic error")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

package storage

import (
	"context"

	"github.com/singh-anurag-7991/shield/internal/limiter"
)

type Storage interface {
	GetLimiter(ctx context.Context, key string) (limiter.Limiter, error)
	SetLimiter(ctx context.Context, key string, limiter limiter.Limiter) error
}

package limiter

import (
	"context"
	"time"

	"github.com/singh-anurag-7991/shield/internal/models"
	"github.com/singh-anurag-7991/shield/internal/storage"
)

type LimiterFactory struct {
	storage storage.Storage
}

func NewFactory(s storage.Storage) *LimiterFactory {
	return &LimiterFactory{storage: s}
}

func (f *LimiterFactory) Create(cfg models.LimiterConfig) Limiter {
	ctx := context.Background()
	l, err := f.storage.GetLimiter(ctx, cfg.Name)
	if err == nil {
		return l
	}

	// Create new limiter
	var lim Limiter
	switch cfg.Algorithm {
	case "token":
		lim = NewTokenBucket(cfg.Capacity, cfg.Rate)
	case "leaky":
		lim = NewLeakyBucket(cfg.Capacity, cfg.Rate)
	case "fixed":
		window, _ := time.ParseDuration(cfg.Window)
		lim = NewFixedWindow(cfg.Capacity, window)
	case "sliding":
		window, _ := time.ParseDuration(cfg.Window)
		lim = NewSlidingLog(cfg.Capacity, window)
	default:
		panic("unknown algorithm: " + cfg.Algorithm)
	}

	// Cache it
	if err := f.storage.SetLimiter(ctx, cfg.Name, lim); err != nil {
		// In production, log this
	}
	return lim
}

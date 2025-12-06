package limiter

import (
	"context"
	"time"

	"github.com/singh-anurag-7991/shield/internal/models"
	"github.com/singh-anurag-7991/shield/internal/rate"
)

type LimiterFactory struct {
	storage rate.Storage
}

func NewFactory(s rate.Storage) *LimiterFactory {
	return &LimiterFactory{storage: s}
}

func (f *LimiterFactory) Create(cfg models.LimiterConfig) rate.Limiter {
	ctx := context.Background()
	if l, err := f.storage.GetLimiter(ctx, cfg.Name); err == nil {
		return l
	}

	var lim rate.Limiter
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

	_ = f.storage.SetLimiter(ctx, cfg.Name, lim)
	return lim
}

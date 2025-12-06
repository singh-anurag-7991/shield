package rate

import "context"

type Limiter interface {
	Allow(key string) bool
	GetStats(key string) LimiterStats
	LimiterType() string
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

type LimiterStats struct {
	Remaining int64
	Limit     int64
	Reset     int64
}

type Storage interface {
	GetLimiter(ctx context.Context, key string) (Limiter, error)
	SetLimiter(ctx context.Context, key string, l Limiter) error
}

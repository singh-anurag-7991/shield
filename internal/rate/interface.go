package rate

type Limiter interface {
	Allow(key string) bool
	GetStats(key string) LimiterStats
}

type LimiterStats struct {
	Remaining int64
	Limit     int64
	Reset     int64
}

package limiter

import (
	"sync"
	"time"

	"github.com/singh-anurag-7991/shield/internal/rate"
)

type FixedWindow struct {
	mu     sync.Mutex
	counts map[string]int64
	resets map[string]time.Time
	limit  int64
	window time.Duration
}

func NewFixedWindow(limit int64, window time.Duration) *FixedWindow {
	fw := &FixedWindow{
		counts: make(map[string]int64),
		resets: make(map[string]time.Time),
		limit:  limit,
		window: window,
	}
	go fw.cleanup()
	return fw
}

func (fw *FixedWindow) Allow(key string) bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now()
	if reset, ok := fw.resets[key]; ok && now.Before(reset) {
		if fw.counts[key] >= fw.limit {
			return false
		}
		fw.counts[key]++
		return true
	}

	fw.counts[key] = 1
	fw.resets[key] = now.Add(fw.window)
	return true
}

func (fw *FixedWindow) GetStats(key string) rate.LimiterStats {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	count := fw.counts[key]
	remaining := fw.limit - count
	if remaining < 0 {
		remaining = 0
	}

	var reset int64
	if r, ok := fw.resets[key]; ok {
		reset = r.Unix()
	}
	return rate.LimiterStats{Remaining: remaining, Limit: fw.limit, Reset: reset}
}

func (fw *FixedWindow) cleanup() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		fw.mu.Lock()
		now := time.Now()
		for k, t := range fw.resets {
			if now.After(t) {
				delete(fw.counts, k)
				delete(fw.resets, k)
			}
		}
		fw.mu.Unlock()
	}
}

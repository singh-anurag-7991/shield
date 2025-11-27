package limiter

import (
	"container/list"
	"sync"
	"time"
)

type SlidingLog struct {
	limit  int64
	window time.Duration
	logs   map[string]*list.List
	mu     sync.Mutex
}

func NewSlidingLog(limit int64, window time.Duration) *SlidingLog {
	return &SlidingLog{
		limit:  limit,
		window: window,
		logs:   make(map[string]*list.List),
	}
}

func (sl *SlidingLog) Allow(key string) bool {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-sl.window)
	if sl.logs[key] == nil {
		sl.logs[key] = list.New()
	}
	l := sl.logs[key]

	// Remove old timestamps
	for e := l.Front(); e != nil; {
		if e.Value.(time.Time).Before(cutoff) {
			next := e.Next()
			l.Remove(e)
			e = next
		} else {
			break
		}
	}

	if int64(l.Len()) < sl.limit {
		l.PushBack(now)
		return true
	}
	return false
}

func (sl *SlidingLog) GetStats(key string) LimiterStats {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	if sl.logs[key] == nil {
		return LimiterStats{Remaining: sl.limit, Limit: sl.limit, Reset: 0}
	}
	count := int64(sl.logs[key].Len())
	return LimiterStats{
		Remaining: max(0, sl.limit-count),
		Limit:     sl.limit,
		Reset:     0, // Sliding window doesn't have a fixed reset
	}
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

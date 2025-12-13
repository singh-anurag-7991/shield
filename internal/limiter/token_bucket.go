package limiter

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/singh-anurag-7991/shield/internal/rate"
)

type TokenBucket struct {
	mu         sync.Mutex
	buckets    map[string]*tbState
	capacity   int64
	refillRate int64
}

type tbState struct {
	tokens     float64
	lastRefill time.Time
}

func NewTokenBucket(capacity, refillRate int64) *TokenBucket {
	return &TokenBucket{
		buckets:    make(map[string]*tbState),
		capacity:   capacity,
		refillRate: refillRate,
	}
}

func (tb *TokenBucket) Allow(key string) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	state := tb.getOrCreate(key)

	now := time.Now()
	elapsed := now.Sub(state.lastRefill).Seconds()
	state.tokens += elapsed * float64(tb.refillRate)
	if state.tokens > float64(tb.capacity) {
		state.tokens = float64(tb.capacity)
	}
	state.lastRefill = now

	if state.tokens >= 1 {
		state.tokens -= 1
		return true
	}
	return false
}

func (tb *TokenBucket) GetStats(key string) rate.LimiterStats {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	state := tb.getOrCreate(key)
	remaining := int64(state.tokens)
	if remaining < 0 {
		remaining = 0
	}
	return rate.LimiterStats{
		Remaining: remaining,
		Limit:     tb.capacity,
		Reset:     time.Now().Add(time.Duration((float64(tb.capacity)-state.tokens)/float64(tb.refillRate)) * time.Second).Unix(),
	}
}

func (tb *TokenBucket) getOrCreate(key string) *tbState {
	if s, ok := tb.buckets[key]; ok {
		return s
	}
	s := &tbState{tokens: float64(tb.capacity), lastRefill: time.Now()}
	tb.buckets[key] = s
	return s
}

func (tb *TokenBucket) LimiterType() string {
	return "token"
}

func (tb *TokenBucket) MarshalJSON() ([]byte, error) {
	type Alias TokenBucket
	return json.Marshal(&struct {
		*Alias
	}{Alias: (*Alias)(tb)})
}

func (tb *TokenBucket) UnmarshalJSON(data []byte) error {
	type Alias TokenBucket
	var a Alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	*tb = TokenBucket(a)
	if tb.buckets == nil {
		tb.buckets = make(map[string]*tbState)
	}
	return nil
}

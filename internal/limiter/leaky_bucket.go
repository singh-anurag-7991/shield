package limiter

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/singh-anurag-7991/shield/internal/rate"
)

type LeakyBucket struct {
	mu       sync.Mutex
	buckets  map[string]*lbState
	capacity int64
	rate     int64
}

type lbState struct {
	water    float64
	lastLeak time.Time
}

func NewLeakyBucket(capacity, rate int64) *LeakyBucket {
	return &LeakyBucket{
		buckets:  make(map[string]*lbState),
		capacity: capacity,
		rate:     rate,
	}
}

func (lb *LeakyBucket) Allow(key string) bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	state := lb.getOrCreate(key)
	now := time.Now()
	elapsed := now.Sub(state.lastLeak).Seconds()
	leaked := elapsed * float64(lb.rate)
	state.water = maxFloat(0.0, state.water-leaked)
	state.lastLeak = now

	if state.water < float64(lb.capacity) {
		state.water += 1
		return true
	}
	return false
}

func (lb *LeakyBucket) GetStats(key string) rate.LimiterStats {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	state := lb.getOrCreate(key)
	waterInt := int64(state.water)
	remaining := lb.capacity - waterInt
	if remaining < 0 {
		remaining = 0
	}
	return rate.LimiterStats{
		Remaining: remaining,
		Limit:     lb.capacity,
		Reset:     0,
	}
}

func (lb *LeakyBucket) getOrCreate(key string) *lbState {
	if lb.buckets == nil { // ← YE SAFETY CHECK ADD KARO
		lb.buckets = make(map[string]*lbState)
	}
	if s, ok := lb.buckets[key]; ok {
		return s
	}
	s := &lbState{water: 0, lastLeak: time.Now()}
	lb.buckets[key] = s
	return s
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func (lb *LeakyBucket) LimiterType() string {
	return "leaky"
}

func (lb *LeakyBucket) MarshalJSON() ([]byte, error) {
	type Alias LeakyBucket
	return json.Marshal(&struct{ *Alias }{Alias: (*Alias)(lb)})
}

func (lb *LeakyBucket) UnmarshalJSON(data []byte) error {
	type Alias LeakyBucket
	var a Alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	*lb = LeakyBucket(a)
	if lb.buckets == nil { // ← SAFETY CHECK AFTER UNMARSHAL
		lb.buckets = make(map[string]*lbState)
	}
	return nil
}

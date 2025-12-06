// internal/limiter/types/redis.go
package types

import (
	"github.com/singh-anurag-7991/shield/internal/limiter"
	"github.com/singh-anurag-7991/shield/internal/rate"
)

// Yeh sirf type aliases hain — no import cycle
func NewTokenBucket() rate.Limiter { return &limiter.TokenBucket{} }
func NewLeakyBucket() rate.Limiter { return &limiter.LeakyBucket{} }
func NewFixedWindow() rate.Limiter { return &limiter.FixedWindow{} }
func NewSlidingLog() rate.Limiter  { return &limiter.SlidingLog{} }

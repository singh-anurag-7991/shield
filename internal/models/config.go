package models

type RateLimitConfig struct {
	Limiters []LimiterConfig `json:"limiters"`
}

type LimiterConfig struct {
	Name      string `json:"name"`
	Algorithm string `json:"algorithm"`
	Capacity  int64  `json:"capacity"`
	Rate      int64  `json:"rate,omitempty"`
	Window    string `json:"window,omitempty"` // e.g., "60s" for fixed/sliding
}

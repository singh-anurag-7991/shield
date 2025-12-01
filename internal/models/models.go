package models

type LimiterConfig struct {
	Name      string `json:"name"`
	Algorithm string `json:"algorithm"`
	Capacity  int64  `json:"capacity"`
	Rate      int64  `json:"rate,omitempty"`
	Window    string `json:"window,omitempty"`
}

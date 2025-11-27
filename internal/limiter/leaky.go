package limiter

// type LeakyBucket struct {
// 	mu       sync.Mutex
// 	buckets  map[string]*leakyState
// 	capacity int64
// 	rate     int64 // leaks per second
// }

// type leakyState struct {
// 	water    float64
// 	lastLeak time.Time
// }

// Same pattern: get or create state for key, leak water, etc.

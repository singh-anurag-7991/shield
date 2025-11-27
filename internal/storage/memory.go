package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/singh-anurag-7991/shield/internal/limiter"
)

type MemoryStorage struct {
	limiters map[string]limiter.Limiter
	mu       sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		limiters: make(map[string]limiter.Limiter),
	}
}

func (ms *MemoryStorage) GetLimiter(ctx context.Context, key string) (limiter.Limiter, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if l, ok := ms.limiters[key]; ok {
		return l, nil
	}
	return nil, fmt.Errorf("limiter not found for key: %s", key)
}

func (ms *MemoryStorage) SetLimiter(ctx context.Context, key string, l limiter.Limiter) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.limiters[key] = l
	return nil
}

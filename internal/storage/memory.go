package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/singh-anurag-7991/shield/internal/rate"
)

type Storage interface {
	GetLimiter(ctx context.Context, key string) (rate.Limiter, error)
	SetLimiter(ctx context.Context, key string, limiter rate.Limiter) error
}

type MemoryStorage struct {
	limiters map[string]rate.Limiter
	mu       sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		limiters: make(map[string]rate.Limiter),
	}
}

func (m *MemoryStorage) GetLimiter(ctx context.Context, key string) (rate.Limiter, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if l, ok := m.limiters[key]; ok {
		return l, nil
	}
	return nil, fmt.Errorf("limiter not found: %s", key)
}

func (m *MemoryStorage) SetLimiter(ctx context.Context, key string, limiter rate.Limiter) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.limiters[key] = limiter
	return nil
}

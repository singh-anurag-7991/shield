package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/singh-anurag-7991/shield/internal/limiter"
	"github.com/singh-anurag-7991/shield/internal/rate"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(url string) (*RedisStorage, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}
	return &RedisStorage{client: client}, nil
}

func (r *RedisStorage) GetLimiter(ctx context.Context, key string) (rate.Limiter, error) {
	data, err := r.client.Get(ctx, "shield:limiter:"+key).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("limiter not found")
	}
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}

	var lim rate.Limiter
	switch wrapper.Type {
	case "token":
		lim = &limiter.TokenBucket{}
	case "leaky":
		lim = &limiter.LeakyBucket{}
	case "fixed":
		lim = &limiter.FixedWindow{}
	case "sliding":
		lim = &limiter.SlidingLog{}
	default:
		return nil, fmt.Errorf("unknown limiter type: %s", wrapper.Type)
	}

	if err := lim.UnmarshalJSON(wrapper.Data); err != nil {
		return nil, err
	}
	return lim, nil
}

func (r *RedisStorage) SetLimiter(ctx context.Context, key string, l rate.Limiter) error {
	data, err := l.MarshalJSON()
	if err != nil {
		return err
	}

	wrapper := struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}{
		Type: l.LimiterType(),
		Data: data,
	}

	payload, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, "shield:limiter:"+key, payload, 24*time.Hour).Err()
}

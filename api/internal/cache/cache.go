package cache

import (
	"context"
	"fmt"
	"sync"
)

var _ Cache = (*InMemoryCache)(nil)

type Cache interface {
	Store(ctx context.Context, key string, value any, ttl int64) error
	Get(ctx context.Context, key string) (interface{}, error)
}

type InMemoryCache struct {
	storage sync.Map
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{}
}

func (cache *InMemoryCache) Get(_ context.Context, key string) (any, error) {
	v, ok := cache.storage.Load(key)
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return v, nil
}

func (cache *InMemoryCache) Store(_ context.Context, key string, value any, ttl int64) error {
	cache.storage.Store(key, value)
	return nil
}

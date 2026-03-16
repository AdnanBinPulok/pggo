package cache

import (
	"sync"
	"time"
)

type item[T any] struct {
	value     T
	expiresAt time.Time
}

// MemoryCache is a simple TTL cache for typed table rows.
type MemoryCache[T any] struct {
	mu    sync.RWMutex
	ttl   time.Duration
	items map[string]item[T]
}

// NewMemoryCache creates a typed memory cache with TTL.
func NewMemoryCache[T any](ttl time.Duration) *MemoryCache[T] {
	return &MemoryCache[T]{ttl: ttl, items: map[string]item[T]{}}
}

// Set stores one key-value item.
func (c *MemoryCache[T]) Set(key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = item[T]{value: value, expiresAt: time.Now().Add(c.ttl)}
}

// Get returns one key-value item if present and not expired.
func (c *MemoryCache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	obj, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		var zero T
		return zero, false
	}
	if !obj.expiresAt.IsZero() && time.Now().After(obj.expiresAt) {
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		var zero T
		return zero, false
	}
	return obj.value, true
}

// Delete removes one cache key.
func (c *MemoryCache[T]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear removes all cache keys.
func (c *MemoryCache[T]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = map[string]item[T]{}
}

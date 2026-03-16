package cache

import "time"

// Manager wraps typed cache operations.
type Manager[T any] struct {
	store *MemoryCache[T]
}

// NewManager creates a cache manager.
func NewManager[T any](ttl time.Duration) *Manager[T] {
	return &Manager[T]{store: NewMemoryCache[T](ttl)}
}

// Set writes one cache item.
func (m *Manager[T]) Set(key string, value T) { m.store.Set(key, value) }

// Get reads one cache item.
func (m *Manager[T]) Get(key string) (T, bool) { return m.store.Get(key) }

// Delete removes one cache item.
func (m *Manager[T]) Delete(key string) { m.store.Delete(key) }

// Clear clears cache content.
func (m *Manager[T]) Clear() { m.store.Clear() }

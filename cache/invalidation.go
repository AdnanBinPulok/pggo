package cache

// InvalidateMany deletes all provided keys.
func (m *Manager[T]) InvalidateMany(keys []string) {
	for _, k := range keys {
		m.Delete(k)
	}
}

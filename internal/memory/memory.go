package memory

import (
	"context"
	"sync"
)

// InMemoryStorage implements a simple in-memory storage system
type InMemoryStorage struct {
	mu      sync.RWMutex
	data    []byte
}

// NewInMemoryStorage creates a new instance of InMemoryStorage
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{}
}

// LoadMemory retrieves the stored memory data
func (m *InMemoryStorage) LoadMemory(ctx context.Context) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if m.data == nil {
			return []byte{}, nil
		}
		return m.data, nil
	}
}

// SaveMemory stores new memory data
func (m *InMemoryStorage) SaveMemory(ctx context.Context, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		m.data = make([]byte, len(data))
		copy(m.data, data)
		return nil
	}
}

// Clear removes all stored memory data
func (m *InMemoryStorage) Clear(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		m.data = nil
		return nil
	}
} 
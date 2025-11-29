// Package cache provides a simple in-memory cache with TTL support
package cache

import (
	"sync"
	"time"
)

// Entry represents a cached value with expiration
type Entry struct {
	Value      interface{}
	Expiration time.Time
}

// Cache is a thread-safe in-memory cache with TTL
type Cache struct {
	data map[string]Entry
	mu   sync.RWMutex
	ttl  time.Duration
}

// New creates a new cache with the specified TTL
func New(ttl time.Duration) *Cache {
	c := &Cache{
		data: make(map[string]Entry),
		ttl:  ttl,
	}

	// Start cleanup goroutine
	go c.cleanupExpired()

	return c
}

// Get retrieves a value from the cache
// Returns the value and true if found and not expired, nil and false otherwise
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.Expiration) {
		return nil, false
	}

	return entry.Value, true
}

// Set stores a value in the cache with the default TTL
func (c *Cache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.ttl)
}

// SetWithTTL stores a value in the cache with a custom TTL
func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = Entry{
		Value:      value,
		Expiration: time.Now().Add(ttl),
	}
}

// Delete removes a value from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

// Clear removes all values from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]Entry)
}

// Size returns the number of items in the cache (including expired)
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.data)
}

// cleanupExpired periodically removes expired entries
func (c *Cache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.data {
			if now.After(entry.Expiration) {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}

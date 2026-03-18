// Package cache provides a thread-safe schema cache for the swagify package.
package cache

import "sync"

// SchemaCache is a thread-safe cache for generated JSON schemas.
// It maps type names to their generated schema representations.
type SchemaCache struct {
	mu    sync.RWMutex
	store map[string]any
}

// New creates a new SchemaCache.
func New() *SchemaCache {
	return &SchemaCache{
		store: make(map[string]any),
	}
}

// Get retrieves a cached schema by type name. Returns the value and whether it was found.
func (c *SchemaCache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.store[key]
	return v, ok
}

// Set stores a schema in the cache by type name.
func (c *SchemaCache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
}

// Has checks if a schema exists in the cache.
func (c *SchemaCache) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.store[key]
	return ok
}

// Keys returns all cached type names.
func (c *SchemaCache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]string, 0, len(c.store))
	for k := range c.store {
		keys = append(keys, k)
	}
	return keys
}

// All returns a copy of all cached schemas.
func (c *SchemaCache) All() map[string]any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make(map[string]any, len(c.store))
	for k, v := range c.store {
		result[k] = v
	}
	return result
}

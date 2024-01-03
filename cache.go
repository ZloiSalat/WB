package main

import (
	"sync"
	"time"
)

type InMemoryCache interface {
	Set(string, interface{}, time.Duration)
	Get(string) (interface{}, bool)
	Remove(string)
	IsEmpty() bool
	GetAll() map[string]interface{}
}

// Cache is a simple in-memory cache.
type Cache struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

// NewCache creates a new instance of the Cache.
func NewCache() (*Cache, error) {
	return &Cache{
		data: make(map[string]interface{}),
	}, nil
}

// Set adds a value to the cache with a specified key and expiration time.
func (c *Cache) Set(key string, value interface{}, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value

	// Schedule automatic removal after expiration
	if expiration > 0 {
		time.AfterFunc(expiration, func() {
			c.Remove(key)
		})
	}
}

// Get retrieves a value from the cache based on the key.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.data[key]
	return value, ok
}

// Remove removes a key-value pair from the cache.
func (c *Cache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

func (c *Cache) IsEmpty() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.data) == 0

}

func (c *Cache) GetAll() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]interface{})
	for key, value := range c.data {
		result[key] = value
	}
	return result
}

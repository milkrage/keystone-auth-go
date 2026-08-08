package cache

import (
	"sync"
	"time"
)

type Cache struct {
	value     string
	expiresAt time.Time
	mu        sync.RWMutex
}

func (c *Cache) Get() (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.value != "" && time.Until(c.expiresAt) > 0 {
		return c.value, true
	}

	return "", false
}

func (c *Cache) Set(value string, expiresAt time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.value = value
	c.expiresAt = expiresAt
}

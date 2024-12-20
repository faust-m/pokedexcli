package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := Cache{
		entries: map[string]cacheEntry{},
		mu:      sync.Mutex{},
	}
	go c.reapLoop(interval)
	return &c
}

func (c *Cache) Add(key string, val []byte) {
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       []byte{},
	}
	entry.val = append(entry.val, val...)
	c.mu.Lock()
	c.entries[key] = entry
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	entry, found := c.entries[key]
	c.mu.Unlock()
	if !found {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for t := range ticker.C {
		for k, v := range c.entries {
			if t.After(v.createdAt.Add(interval)) {
				c.mu.Lock()
				delete(c.entries, k)
				c.mu.Unlock()
			}
		}
	}
}

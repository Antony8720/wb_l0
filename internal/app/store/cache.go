package store

import (
	"sync"
)

type Cache struct {
	sync.RWMutex
	data map[string][]byte
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string][]byte),
	}
}

func (cache *Cache) Get(id string) ([]byte, bool) {
	cache.RLock()
	defer cache.RUnlock()
	value, ok := cache.data[id]
	return value, ok
}

func (cache *Cache) Set(id string, value []byte) {
	cache.Lock()
	defer cache.Unlock()
	cache.data[id] = value
}

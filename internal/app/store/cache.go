package store

import (
	"log"
	"sync"
)

type Cache struct {
	sync.RWMutex
	data map[string][]byte
}

func NewCache() *Cache{
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
	log.Printf("cache: %v", cache.data[id])
	cache.data[id] = value
	log.Printf("id: %v", id)
	log.Printf("cache: %v", cache.data[id])
}
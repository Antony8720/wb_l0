package store

import "sync"

type Cache struct {
	sync.RWMutex
	data map[int]string
}

func NewCache() *Cache{
	return &Cache{
		data: make(map[int]string),
	}
}

func (cache *Cache) Get(id int) (string, bool) {
	cache.RLock()
	defer cache.RUnlock()
	value, ok := cache.data[id]
	return value, ok
}

func (cache *Cache) Set(id int, value string) {
	cache.Lock()
	defer cache.Unlock()
	cache.data[id] = value
}
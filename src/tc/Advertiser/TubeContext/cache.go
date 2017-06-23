package main

import (
	"sync"
	"time"
)

type (
	Cache struct {
		sync.RWMutex
		items    cacheItems
		duration time.Duration
	}
	cacheItem struct {
		expires time.Time
		data    interface{}
	}
	cacheItems map[string]*cacheItem
)

func NewCache(duration time.Duration) *Cache {
	return &Cache{
		duration: duration,
		items:    make(cacheItems),
	}
}

func (c *Cache) Put(key string, data interface{}) {
	c.Lock()
	defer c.Unlock()

	c.items[key] = &cacheItem{
		expires: time.Now().Add(c.duration),
		data:    data,
	}
}

func (c *Cache) Get(key string, oName string) (data interface{}, exists bool, expired bool) {
	defer do(measure("cache-get " + oName))
	c.RLock()
	defer c.RUnlock()

	var e *cacheItem

	if e, exists = c.items[key]; exists {
		data = e.data
		expired = e.expires.Before(time.Now())
	}

	return data, exists, expired
}

func (c *Cache) Delete(key string) {
	defer do(measure("cache-delete"))
	c.Lock()
	defer c.Unlock()
	delete(c.items, key)
}

func (c *Cache) Len() int {
	return len(c.items)
}

package global

import "sync"

type RowCache struct {
	lock  sync.RWMutex
	cache map[uint32]bool
}

func NewRowCache() *RowCache {
	rowCache := &RowCache{}
	rowCache.cache = make(map[uint32]bool)
	return rowCache
}

func (c *RowCache) Put(key []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()
	hash := Hash32(key)
	c.cache[hash] = true
}

func (c *RowCache) Get(key []byte) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	hash := Hash32(key)
	_, ok := c.cache[hash]
	return ok
}

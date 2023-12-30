package cache

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

type value struct {
	v          []byte
	expiration int64
}

type Cache struct {
	sync.RWMutex

	values                 map[string]*value
	garbageCollectInterval time.Duration
}

func New(ctx context.Context, garbageCollectInterval time.Duration) *Cache {
	c := &Cache{
		values:                 make(map[string]*value),
		garbageCollectInterval: garbageCollectInterval,
	}

	go c.deleteOldValues(ctx)

	return c
}

func (c *Cache) Set(key string, val interface{}, timeout time.Duration) {
	var expiration int64
	if timeout != 0 {
		expiration = time.Now().Add(timeout).UnixNano()
	}

	data, _ := json.Marshal(val)

	// lock mutex
	c.Lock()

	c.values[key] = &value{
		v:          data,
		expiration: expiration,
	}

	// unlock mutex
	c.Unlock()
}

func (c *Cache) Get(key string, dst interface{}) bool {
	c.RLock()

	val, ok := c.values[key]

	c.RUnlock()

	if !ok {
		return false
	}

	err := json.Unmarshal(val.v, &dst)
	if err != nil {
		return false
	}

	return true
}

func (c *Cache) Size() int {
	c.RLock()
	size := len(c.values)
	c.RUnlock()

	return size
}

func (c *Cache) deleteOldValues(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(c.garbageCollectInterval):
			c.Lock()

			now := time.Now().UnixNano()

			for key, val := range c.values {
				if val.expiration > 0 && now > val.expiration {
					delete(c.values, key)
				}
			}

			c.Unlock()
		}
	}
}

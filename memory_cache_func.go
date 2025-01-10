package memory_cache

import "errors"

func New() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	c.data[key] = value
}

func (c *Cache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	value, ok := c.data[key]
	if !ok {
		c.mu.RUnlock()
		return nil, errors.New("key not found")
	}
	c.mu.RUnlock()
	return value, nil
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
}

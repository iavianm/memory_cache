package memorycache

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrKeyExpired  = errors.New("key expired")
)

type entry[T any] struct {
	value     T
	expiresAt time.Time
}

func (e entry[T]) isExpired() bool {
	if e.expiresAt.IsZero() {
		return false
	}
	return time.Now().After(e.expiresAt)
}

func NewMemoryCache[T any](defaultTTL, cleanupInterval time.Duration) *MemoryCache[T] {
	mc := &MemoryCache[T]{
		cache:      make(map[string]entry[T]),
		mu:         sync.RWMutex{},
		defaultTTL: defaultTTL,
		stop:       make(chan struct{}),
	}

	if cleanupInterval > 0 {
		go mc.startJanitor(cleanupInterval)
	}

	return mc
}

func (mc *MemoryCache[T]) Set(key string, value T) {
	mc.SetWithTTL(key, value, mc.defaultTTL)
}

func (mc *MemoryCache[T]) SetWithTTL(key string, value T, ttl time.Duration) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}
	mc.cache[key] = entry[T]{
		value:     value,
		expiresAt: expiresAt,
	}
}

func (mc *MemoryCache[T]) Get(key string) (T, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	var zero T
	e, ok := mc.cache[key]

	if !ok {
		return zero, ErrKeyNotFound
	}

	if e.isExpired() {
		return zero, ErrKeyExpired
	}

	return e.value, nil
}

func (mc *MemoryCache[T]) Delete(key string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	delete(mc.cache, key)
}

func (mc *MemoryCache[T]) deleteExpired() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	for k, v := range mc.cache {
		if v.isExpired() {
			delete(mc.cache, k)
		}
	}
}

func (mc *MemoryCache[T]) startJanitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mc.deleteExpired()
		case <-mc.stop:
			return
		}
	}
}

func (mc *MemoryCache[T]) Close() {
	mc.closeOnce.Do(func() {
		close(mc.stop)
	})
}

func (mc *MemoryCache[T]) size() int {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return len(mc.cache)
}

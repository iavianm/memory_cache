package memorycache

import (
	"sync"
	"time"
)

type MemoryCache[T any] struct {
	cache      map[string]entry[T]
	mu         sync.RWMutex
	defaultTTL time.Duration
	stop       chan struct{}
	closeOnce  sync.Once
}

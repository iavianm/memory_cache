package memory_cache

import "sync"

type Cache struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

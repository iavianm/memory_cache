package memorycache

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMemoryCache_Set_Get(t *testing.T) {
	cache := NewMemoryCache[string](time.Minute, 0)
	defer cache.Close()
	cache.Set("foo", "bar")

	got, err := cache.Get("foo")
	if err != nil {
		t.Fatal(err)
	}
	if got != "bar" {
		t.Errorf("got %q, want %q", got, "bar")
	}
}

func TestMemoryCache_Set_Expire(t *testing.T) {
	cache := NewMemoryCache[string](0, 0)
	defer cache.Close()

	cache.SetWithTTL("foo", "bar", 500*time.Millisecond)
	_, err := cache.Get("foo")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(800 * time.Millisecond)

	_, err = cache.Get("foo")
	if !errors.Is(err, ErrKeyExpired) {
		t.Errorf("got %v, want %v", err, ErrKeyExpired)
	}
}

func TestMemoryCache_Set_Delete(t *testing.T) {
	cache := NewMemoryCache[string](time.Minute, 0)
	defer cache.Close()

	cache.Set("foo", "bar")
	if _, err := cache.Get("foo"); err != nil {
		t.Fatal(err)
	}

	cache.Delete("foo")
	_, err := cache.Get("foo")
	if !errors.Is(err, ErrKeyNotFound) {
		t.Errorf("got %v, want %v", err, ErrKeyNotFound)
	}
}

func TestMemoryCache_Janitor(t *testing.T) {
	cache := NewMemoryCache[string](0, 20*time.Millisecond)
	defer cache.Close()

	cache.SetWithTTL("foo", "bar", 10*time.Millisecond)
	if n := cache.size(); n != 1 {
		t.Fatalf("got %d, want %d", n, 1)
	}
	time.Sleep(100 * time.Millisecond)
	if n := cache.size(); n != 0 {
		t.Errorf("got %d, want %d", n, 0)
	}
}

func TestMemoryCache_Concurrent(t *testing.T) {
	cache := NewMemoryCache[int](time.Minute, 10*time.Millisecond)
	defer cache.Close()

	const goroutines = 50

	var wg sync.WaitGroup
	for i := range goroutines {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("key-%d", i)
			cache.Set(key, i)
			_, _ = cache.Get(key)
			_, _ = cache.Get("key-0")
			cache.Delete(key)
		}(i)
	}
	wg.Wait()
}

func TestMemoryCache_Set_NoExpire(t *testing.T) {
	cache := NewMemoryCache[string](0, 0)
	defer cache.Close()

	cache.SetWithTTL("foo", "bar", 0)

	time.Sleep(100 * time.Millisecond)

	got, err := cache.Get("foo")
	if err != nil {
		t.Fatal(err)
	}
	if got != "bar" {
		t.Errorf("got %q, want %q", got, "bar")
	}
}

# memory_cache

[English](#english) | [Русский](#русский)

---

## English

A generic thread-safe in-memory cache for Go with TTL support and background
cleanup of expired entries.

- generics — no `interface{}`, no type assertions
- TTL per cache (default) or per key
- background janitor removes expired entries
- safe for concurrent use (`sync.RWMutex`)

### Installation

```bash
go get github.com/iavianm/memory_cache
```

### API

| Method | Description |
| --- | --- |
| `NewMemoryCache[T](defaultTTL, cleanupInterval time.Duration) *MemoryCache[T]` | Creates a cache. `defaultTTL` is used by `Set`; `0` means entries never expire. `cleanupInterval` is the janitor period; `0` disables background cleanup. |
| `Set(key string, value T)` | Stores a value with `defaultTTL`. |
| `SetWithTTL(key string, value T, ttl time.Duration)` | Stores a value with an explicit TTL; `0` means it never expires. |
| `Get(key string) (T, error)` | Returns the value, or the zero value of `T` plus an error. |
| `Delete(key string)` | Removes the key. |
| `Close()` | Stops the background janitor. Safe to call multiple times and from multiple goroutines. |

Errors returned by `Get` (compare with `errors.Is`):

| Error | Meaning |
| --- | --- |
| `ErrKeyNotFound` | The key is not in the cache. |
| `ErrKeyExpired` | The key exists but its TTL has passed. |

### Usage

```go
package main

import (
	"errors"
	"fmt"
	"time"

	memorycache "github.com/iavianm/memory_cache"
)

func main() {
	// default TTL 5 minutes, cleanup every minute
	cache := memorycache.NewMemoryCache[string](5*time.Minute, time.Minute)
	defer cache.Close()

	cache.Set("user", "Alice")                       // lives 5 minutes
	cache.SetWithTTL("token", "abc", 30*time.Second) // lives 30 seconds
	cache.SetWithTTL("config", "prod", 0)            // lives forever

	value, err := cache.Get("user")
	switch {
	case errors.Is(err, memorycache.ErrKeyNotFound):
		fmt.Println("no such key")
	case errors.Is(err, memorycache.ErrKeyExpired):
		fmt.Println("expired")
	case err != nil:
		fmt.Println("error:", err)
	default:
		fmt.Println(value) // Alice
	}

	cache.Delete("user")
}
```

### Notes

- `defer cache.Close()` right after creation. The janitor goroutine holds a
  reference to the cache, so without `Close` it is never garbage collected.
- An expired entry is not returned by `Get` even before the janitor removes it.
- With `cleanupInterval = 0` expired entries stay in memory until they are
  overwritten or deleted.
- The cache has no size limit and no eviction policy.

### Tests

```bash
go test -race ./...
```

### Requirements

Go 1.26.4 or newer.

---

## Русский

Потокобезопасный обобщённый (generic) in-memory кэш для Go с поддержкой TTL и
фоновой очисткой протухших записей.

- дженерики — никакого `interface{}` и приведения типов
- TTL на весь кэш (по умолчанию) или на отдельный ключ
- фоновая горутина удаляет протухшие записи
- безопасен при конкурентном доступе (`sync.RWMutex`)

### Установка

```bash
go get github.com/iavianm/memory_cache
```

### API

| Метод | Описание |
| --- | --- |
| `NewMemoryCache[T](defaultTTL, cleanupInterval time.Duration) *MemoryCache[T]` | Создаёт кэш. `defaultTTL` использует `Set`; `0` — записи не истекают. `cleanupInterval` — период фоновой очистки; `0` отключает её. |
| `Set(key string, value T)` | Сохраняет значение с `defaultTTL`. |
| `SetWithTTL(key string, value T, ttl time.Duration)` | Сохраняет значение с явным TTL; `0` — вечно. |
| `Get(key string) (T, error)` | Возвращает значение либо нулевое значение `T` и ошибку. |
| `Delete(key string)` | Удаляет ключ. |
| `Close()` | Останавливает фоновую очистку. Можно вызывать повторно и из нескольких горутин. |

Ошибки `Get` (сравнивать через `errors.Is`):

| Ошибка | Значение |
| --- | --- |
| `ErrKeyNotFound` | Ключа нет в кэше. |
| `ErrKeyExpired` | Ключ есть, но его TTL истёк. |

### Пример использования

```go
package main

import (
	"errors"
	"fmt"
	"time"

	memorycache "github.com/iavianm/memory_cache"
)

func main() {
	// TTL по умолчанию 5 минут, очистка раз в минуту
	cache := memorycache.NewMemoryCache[string](5*time.Minute, time.Minute)
	defer cache.Close()

	cache.Set("user", "Alice")                       // живёт 5 минут
	cache.SetWithTTL("token", "abc", 30*time.Second) // живёт 30 секунд
	cache.SetWithTTL("config", "prod", 0)            // живёт вечно

	value, err := cache.Get("user")
	switch {
	case errors.Is(err, memorycache.ErrKeyNotFound):
		fmt.Println("ключа нет")
	case errors.Is(err, memorycache.ErrKeyExpired):
		fmt.Println("протухло")
	case err != nil:
		fmt.Println("ошибка:", err)
	default:
		fmt.Println(value) // Alice
	}

	cache.Delete("user")
}
```

### Примечания

- Ставьте `defer cache.Close()` сразу после создания. Фоновая горутина держит
  ссылку на кэш, поэтому без `Close` сборщик мусора его не заберёт.
- Протухшая запись не возвращается из `Get` ещё до того, как её удалит фоновая
  очистка.
- При `cleanupInterval = 0` протухшие записи остаются в памяти, пока их не
  перезапишут или не удалят.
- У кэша нет ограничения по размеру и политики вытеснения.

### Тесты

```bash
go test -race ./...
```

### Требования

Go 1.26.4 или новее.

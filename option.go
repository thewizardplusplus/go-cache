package cache

import (
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

// Option ...
type Option func(cache *Cache)

// WithStorage ...
//
// Default: an instance of the hashmap.ConcurrentHashMap structure
// with default options.
//
func WithStorage(storage hashmap.Storage) Option {
	return func(cache *Cache) {
		cache.storage = storage
	}
}

// WithClock ...
//
// Default: the time.Now() function.
//
func WithClock(clock Clock) Option {
	return func(cache *Cache) {
		cache.clock = clock
	}
}

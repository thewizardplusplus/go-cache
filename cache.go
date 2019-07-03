package cache

import (
	"time"

	hashmap "github.com/thewizardplusplus/go-hashmap"
)

//go:generate mockery -name=Storage -inpkg -case=underscore -testonly

// Storage ...
type Storage interface {
	Get(key hashmap.Key) (data interface{}, ok bool)
	Set(key hashmap.Key, data interface{})
	Delete(key hashmap.Key)
}

// Clock ...
type Clock func() time.Time

// Cache ...
type Cache struct {
	storage Storage
	clock   Clock
}

type value struct {
	data           interface{}
	expirationTime time.Time
}

// NewCache ...
func NewCache(storage Storage, clock Clock) Cache {
	return Cache{storage, clock}
}

// Set ...
func (cache Cache) Set(key hashmap.Key, data interface{}, ttl time.Duration) {
	var expirationTime time.Time
	if ttl != 0 {
		expirationTime = cache.clock().Add(ttl)
	}

	cache.storage.Set(key, value{data, expirationTime})
}

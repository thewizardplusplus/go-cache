package cache

import (
	"errors"
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

// ...
var (
	ErrKeyMissed  = errors.New("key missed")
	ErrKeyExpired = errors.New("key expired")
)

// NewCache ...
func NewCache(storage Storage, clock Clock) Cache {
	return Cache{storage, clock}
}

// Get ...
func (cache Cache) Get(key hashmap.Key) (data interface{}, err error) {
	data, ok := cache.storage.Get(key)
	if !ok {
		return nil, ErrKeyMissed
	}

	value := data.(value) // nolint: vetshadow
	expirationTime := value.expirationTime
	if !expirationTime.IsZero() && cache.clock().After(expirationTime) {
		return nil, ErrKeyExpired
	}

	return value.data, nil
}

// Set ...
func (cache Cache) Set(key hashmap.Key, data interface{}, ttl time.Duration) {
	var expirationTime time.Time
	if ttl != 0 {
		expirationTime = cache.clock().Add(ttl)
	}

	cache.storage.Set(key, value{data, expirationTime})
}

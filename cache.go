package cache

import (
	"errors"
	"time"

	hashmap "github.com/thewizardplusplus/go-hashmap"
)

// ...
var (
	ErrKeyMissed  = errors.New("key missed")
	ErrKeyExpired = errors.New("key expired")
)

// Clock ...
type Clock func() time.Time

// Cache ...
type Cache struct {
	storage hashmap.Storage
	clock   Clock
}

// NewCache ...
func NewCache(options ...Option) Cache {
	// default options
	cache := Cache{
		storage: hashmap.NewConcurrentHashMap(),
		clock:   time.Now,
	}
	for _, option := range options {
		option(&cache)
	}

	return cache
}

// Get ...
func (cache Cache) Get(key hashmap.Key) (data interface{}, err error) {
	data, ok := cache.storage.Get(key)
	if !ok {
		return nil, ErrKeyMissed
	}

	value := data.(Value)
	if value.IsExpired(cache.clock) {
		return nil, ErrKeyExpired
	}

	return value.Data, nil
}

// GetWithGC ...
func (cache Cache) GetWithGC(key hashmap.Key) (data interface{}, err error) {
	data, err = cache.Get(key)
	if err != nil {
		if err == ErrKeyExpired {
			cache.storage.Delete(key)
		}

		return nil, err
	}

	return data, nil
}

// Set ...
func (cache Cache) Set(key hashmap.Key, data interface{}, ttl time.Duration) {
	var expirationTime time.Time
	if ttl != 0 {
		expirationTime = cache.clock().Add(ttl)
	}

	cache.storage.Set(key, Value{data, expirationTime})
}

// Delete ...
func (cache Cache) Delete(key hashmap.Key) {
	cache.storage.Delete(key)
}

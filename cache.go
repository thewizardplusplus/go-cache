package cache

import (
	"context"
	"errors"
	"time"

	"github.com/thewizardplusplus/go-cache/gc"
	"github.com/thewizardplusplus/go-cache/models"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

// ...
var (
	ErrKeyMissed  = errors.New("key missed")
	ErrKeyExpired = errors.New("key expired")
)

// Cache ...
type Cache struct {
	storage hashmap.Storage
	clock   models.Clock
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

// NewCacheWithGC ...
func NewCacheWithGC(options ...OptionWithGC) Cache {
	// default config
	config := ConfigWithGC{
		ctx:     context.Background(),
		storage: hashmap.NewConcurrentHashMap(),
		clock:   time.Now,
		gcFactory: func(storage hashmap.Storage, clock models.Clock) gc.GC {
			return gc.NewPartialGC(storage, gc.PartialGCWithClock(clock))
		},
		gcPeriod: 100 * time.Millisecond,
	}
	for _, option := range options {
		option(&config)
	}

	gcInstance := config.gcFactory(config.storage, config.clock)
	go gc.Run(config.ctx, gcInstance, config.gcPeriod)

	return NewCache(WithStorage(config.storage), WithClock(config.clock))
}

// Get ...
func (cache Cache) Get(key hashmap.Key) (data interface{}, err error) {
	data, ok := cache.storage.Get(key)
	if !ok {
		return nil, ErrKeyMissed
	}

	value := data.(models.Value)
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

// Iterate ...
func (cache Cache) Iterate(handler hashmap.Handler) bool {
	return cache.storage.Iterate(func(key hashmap.Key, data interface{}) bool {
		value := data.(models.Value)
		if value.IsExpired(cache.clock) {
			return true
		}

		return handler(key, value.Data)
	})
}

// Set ...
func (cache Cache) Set(key hashmap.Key, data interface{}, ttl time.Duration) {
	var expirationTime time.Time
	if ttl != 0 {
		expirationTime = cache.clock().Add(ttl)
	}

	cache.storage.Set(key, models.Value{
		Data:           data,
		ExpirationTime: expirationTime,
	})
}

// Delete ...
func (cache Cache) Delete(key hashmap.Key) {
	cache.storage.Delete(key)
}

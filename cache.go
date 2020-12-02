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
//
// It additionally runs garbage collection in background.
//
func NewCacheWithGC(options ...OptionWithGC) Cache {
	config := newConfigWithGC(options)

	gcInstance := config.gcFactory(config.storage, config.clock)
	go gc.Run(config.ctx, gcInstance, config.gcPeriod)

	return NewCache(WithStorage(config.storage), WithClock(config.clock))
}

// Get ...
//
// The error can be ErrKeyMissed or ErrKeyExpired only.
//
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
//
// It additionally deletes the value if its time to live expired.
//
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
//
// If the handler returns false, iteration is broken.
//
func (cache Cache) Iterate(ctx context.Context, handler hashmap.Handler) bool {
	return cache.iterateWithExpiredHandler(ctx, handler, func(key hashmap.Key) {})
}

// IterateWithGC ...
//
// It additionally deletes iterated values if their time to live expired.
//
// If the handler returns false, iteration is broken.
//
func (cache Cache) IterateWithGC(
	ctx context.Context,
	handler hashmap.Handler,
) bool {
	return cache.iterateWithExpiredHandler(ctx, handler, cache.storage.Delete)
}

// Set ...
//
// Zero time to live means infinite one.
//
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

func (cache Cache) iterateWithExpiredHandler(
	ctx context.Context,
	handler hashmap.Handler,
	expiredHandler func(key hashmap.Key),
) bool {
	return cache.storage.Iterate(
		hashmap.WithInterruption(ctx, func(key hashmap.Key, data interface{}) bool {
			value := data.(models.Value)
			if value.IsExpired(cache.clock) {
				expiredHandler(key)

				return true
			}

			return handler(key, value.Data)
		}),
	)
}

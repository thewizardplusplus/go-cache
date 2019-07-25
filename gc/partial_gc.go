package gc

import (
	cache "github.com/thewizardplusplus/go-cache"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

// PartialGC ...
type PartialGC struct {
	storage Storage
	clock   cache.Clock
}

// NewPartialGC ...
func NewPartialGC(storage Storage, clock cache.Clock) PartialGC {
	return PartialGC{storage, clock}
}

// Clean ...
func (gc PartialGC) Clean() {
	// algorithm is based on expiration in Redis
	// see for details: https://redis.io/commands/expire#how-redis-expires-keys
	for {
		var counter counter // nolint: vetshadow
		gc.storage.Iterate(func(key hashmap.Key, value interface{}) bool {
			if value.(cache.Value).IsExpired(gc.clock) {
				gc.storage.Delete(key)
				counter.expired++
			}

			counter.iterated++
			// iterate over maxIteratedCount values only
			return counter.iterated < maxIteratedCount
		})

		// if a percent of expired values less than minExpiredPercent, stop cleaning
		expiredValuesPercent := float64(counter.expired) / float64(counter.iterated)
		if counter.iterated == 0 || expiredValuesPercent < minExpiredPercent {
			break
		}
	}
}

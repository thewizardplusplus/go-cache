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

const (
	maxIteratedValuesCount  = 20
	minExpiredValuesPercent = 0.25
)

// NewPartialGC ...
func NewPartialGC(storage Storage, clock cache.Clock) PartialGC {
	return PartialGC{storage, clock}
}

// Clean ...
func (gc PartialGC) Clean() {
	// algorithm is based on expiration in Redis
	// see for details: https://redis.io/commands/expire#how-redis-expires-keys
	for {
		var expiredCount int
		var iteratedCount int
		gc.storage.Iterate(func(key hashmap.Key, value interface{}) bool {
			if value.(cache.Value).IsExpired(gc.clock) {
				gc.storage.Delete(key)
				expiredCount++
			}

			iteratedCount++
			// iterate over maxIteratedValuesCount values only
			return iteratedCount < maxIteratedValuesCount
		})

		// if a percent of expired values less than minExpiredValuesPercent,
		// stop cleaning
		expiredValuesPercent := float64(expiredCount) / float64(iteratedCount)
		if iteratedCount == 0 || expiredValuesPercent < minExpiredValuesPercent {
			break
		}
	}
}

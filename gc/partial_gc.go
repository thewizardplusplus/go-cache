package gc

import (
	cache "github.com/thewizardplusplus/go-cache"
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
		iterator := newIterator(gc.storage, gc.clock)
		gc.storage.Iterate(iterator.handleIteration)

		if iterator.stopClean() {
			break
		}
	}
}

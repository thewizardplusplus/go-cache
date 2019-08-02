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
//
// Its algorithm is based on expiration in Redis.
// See for details: https://redis.io/commands/expire#how-redis-expires-keys
func (gc PartialGC) Clean() {
	for {
		iterator := newIterator(gc.storage, gc.clock)
		gc.storage.Iterate(iterator.handleIteration)

		if iterator.stopClean() {
			break
		}
	}
}

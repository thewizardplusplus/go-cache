package gc

import (
	"context"

	cache "github.com/thewizardplusplus/go-cache"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

// PartialGC ...
type PartialGC struct {
	storage           hashmap.Storage
	clock             cache.Clock
	maxIteratedCount  int
	minExpiredPercent float64
}

// NewPartialGC ...
func NewPartialGC(storage hashmap.Storage, clock cache.Clock) PartialGC {
	return PartialGC{storage, clock, maxIteratedCount, minExpiredPercent}
}

// Clean ...
//
// Its algorithm is based on expiration in Redis.
// See for details: https://redis.io/commands/expire#how-redis-expires-keys
func (gc PartialGC) Clean(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			iterator :=
				newIterator(gc.storage, gc.clock, gc.maxIteratedCount, gc.minExpiredPercent)
			gc.storage.Iterate(withInterruption(ctx, iterator.handleIteration))

			if iterator.stopClean() {
				return
			}
		}
	}
}

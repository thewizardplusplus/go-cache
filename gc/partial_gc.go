package gc

import (
	"context"

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
func (gc PartialGC) Clean(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			iterator := newIterator(gc.storage, gc.clock)
			gc.storage.Iterate(withInterruption(ctx, iterator.handleIteration))

			if iterator.stopClean() {
				return
			}
		}
	}
}

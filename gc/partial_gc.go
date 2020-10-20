package gc

import (
	"context"
	"time"

	hashmaputils "github.com/thewizardplusplus/go-cache/hashmap-utils"
	"github.com/thewizardplusplus/go-cache/models"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

// PartialGC ...
type PartialGC struct {
	storage           hashmap.Storage
	clock             models.Clock
	maxIteratedCount  int
	minExpiredPercent float64
}

// NewPartialGC ...
func NewPartialGC(
	storage hashmap.Storage,
	options ...PartialGCOption,
) PartialGC {
	gc := PartialGC{
		storage: storage,

		// default options
		clock:             time.Now,
		maxIteratedCount:  defaultMaxIteratedCount,
		minExpiredPercent: defaultMinExpiredPercent,
	}
	for _, option := range options {
		option(&gc)
	}

	return gc
}

// Clean ...
//
// Its algorithm is based on expiration in Redis. See for details:
// https://redis.io/commands/expire#how-redis-expires-keys
//
func (gc PartialGC) Clean(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			iterator := newIterator(
				gc.storage,
				gc.clock,
				gc.maxIteratedCount,
				gc.minExpiredPercent,
			)
			gc.storage.Iterate(hashmaputils.HandlerWithInterruption(
				ctx,
				iterator.handleIteration,
			))

			if iterator.stopClean() {
				return
			}
		}
	}
}

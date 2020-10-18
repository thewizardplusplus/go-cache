package gc

import (
	"context"

	cache "github.com/thewizardplusplus/go-cache"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

// TotalGC ...
type TotalGC struct {
	storage hashmap.Storage
	clock   cache.Clock
}

// NewTotalGC ...
func NewTotalGC(storage hashmap.Storage, clock cache.Clock) TotalGC {
	return TotalGC{storage, clock}
}

// Clean ...
func (gc TotalGC) Clean(ctx context.Context) {
	gc.storage.Iterate(withInterruption(ctx, gc.handleIteration))
}

func (gc TotalGC) handleIteration(key hashmap.Key, value interface{}) bool {
	if value.(cache.Value).IsExpired(gc.clock) {
		gc.storage.Delete(key)
	}

	return true
}

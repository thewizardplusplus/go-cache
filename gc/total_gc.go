package gc

import (
	"context"
	"time"

	cache "github.com/thewizardplusplus/go-cache"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

// TotalGC ...
type TotalGC struct {
	storage hashmap.Storage
	clock   cache.Clock
}

// NewTotalGC ...
func NewTotalGC(storage hashmap.Storage, options ...TotalGCOption) TotalGC {
	gc := TotalGC{
		storage: storage,
		clock:   time.Now, // default value
	}
	for _, option := range options {
		option(&gc)
	}

	return gc
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

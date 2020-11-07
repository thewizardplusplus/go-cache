package gc

import (
	"context"
	"time"

	"github.com/thewizardplusplus/go-cache/models"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

// TotalGC ...
type TotalGC struct {
	storage hashmap.Storage
	clock   models.Clock
}

// NewTotalGC ...
func NewTotalGC(storage hashmap.Storage, options ...TotalGCOption) TotalGC {
	gc := TotalGC{
		storage: storage,

		// default options
		clock: time.Now,
	}
	for _, option := range options {
		option(&gc)
	}

	return gc
}

// Clean ...
func (gc TotalGC) Clean(ctx context.Context) {
	gc.storage.Iterate(hashmap.WithInterruption(ctx, gc.handleIteration))
}

func (gc TotalGC) handleIteration(key hashmap.Key, value interface{}) bool {
	if value.(models.Value).IsExpired(gc.clock) {
		gc.storage.Delete(key)
	}

	return true
}

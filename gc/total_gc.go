package gc

import (
	"context"

	cache "github.com/thewizardplusplus/go-cache"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

//go:generate mockery -name=Storage -inpkg -case=underscore -testonly

// Storage ...
type Storage interface {
	Iterate(handler hashmap.Handler) bool
	Delete(key hashmap.Key)
}

// TotalGC ...
type TotalGC struct {
	storage Storage
	clock   cache.Clock
}

// NewTotalGC ...
func NewTotalGC(storage Storage, clock cache.Clock) TotalGC {
	return TotalGC{storage, clock}
}

// Clean ...
func (gc TotalGC) Clean(ctx context.Context) {
	gc.storage.Iterate(gc.handleIteration)
}

func (gc TotalGC) handleIteration(key hashmap.Key, value interface{}) bool {
	if value.(cache.Value).IsExpired(gc.clock) {
		gc.storage.Delete(key)
	}

	return true
}

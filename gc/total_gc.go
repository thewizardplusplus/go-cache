package gc

import (
	"context"
	"time"

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
	period  time.Duration
	storage Storage
	clock   cache.Clock
}

// NewTotalGC ...
func NewTotalGC(
	period time.Duration,
	storage Storage,
	clock cache.Clock,
) TotalGC {
	return TotalGC{period, storage, clock}
}

// Clean ...
func (gc TotalGC) Clean() {
	gc.storage.Iterate(func(key hashmap.Key, value interface{}) bool {
		if value.(cache.Value).IsExpired(gc.clock) {
			gc.storage.Delete(key)
		}

		return true
	})
}

// Run ...
func (gc TotalGC) Run(ctx context.Context) {
	ticker := time.NewTicker(gc.period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			gc.Clean()
		case <-ctx.Done():
			return
		}
	}
}

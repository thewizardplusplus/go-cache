package gc

import (
	cache "github.com/thewizardplusplus/go-cache"
)

// PartialGCOption ...
type PartialGCOption func(gc *PartialGC)

// PartialGCWithClock ...
//
// Default: the time.Now() function.
//
func PartialGCWithClock(clock cache.Clock) PartialGCOption {
	return func(gc *PartialGC) {
		gc.clock = clock
	}
}

// PartialGCWithMaxIteratedCount ...
//
// Default: 20.
//
func PartialGCWithMaxIteratedCount(maxIteratedCount int) PartialGCOption {
	return func(gc *PartialGC) {
		gc.maxIteratedCount = maxIteratedCount
	}
}

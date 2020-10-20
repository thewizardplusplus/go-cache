package gc

import (
	"github.com/thewizardplusplus/go-cache/models"
)

const (
	defaultMaxIteratedCount  = 20
	defaultMinExpiredPercent = 0.25
)

// PartialGCOption ...
type PartialGCOption func(gc *PartialGC)

// PartialGCWithClock ...
//
// Default: the time.Now() function.
//
func PartialGCWithClock(clock models.Clock) PartialGCOption {
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

// PartialGCWithMinExpiredPercent ...
//
// Default: 0.25.
//
func PartialGCWithMinExpiredPercent(minExpiredPercent float64) PartialGCOption {
	return func(gc *PartialGC) {
		gc.minExpiredPercent = minExpiredPercent
	}
}

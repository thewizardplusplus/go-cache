package gc

import (
	cache "github.com/thewizardplusplus/go-cache"
)

// TotalGCOption ...
type TotalGCOption func(gc *TotalGC)

// TotalGCWithClock ...
//
// Default: the time.Now() function.
//
func TotalGCWithClock(clock cache.Clock) TotalGCOption {
	return func(gc *TotalGC) {
		gc.clock = clock
	}
}

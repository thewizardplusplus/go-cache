package gc

import (
	"github.com/thewizardplusplus/go-cache/models"
)

// TotalGCOption ...
type TotalGCOption func(gc *TotalGC)

// TotalGCWithClock ...
//
// Default: the time.Now() function.
//
func TotalGCWithClock(clock models.Clock) TotalGCOption {
	return func(gc *TotalGC) {
		gc.clock = clock
	}
}

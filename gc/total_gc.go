package gc

import (
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

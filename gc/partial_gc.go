package gc

import (
	cache "github.com/thewizardplusplus/go-cache"
)

// PartialGC ...
type PartialGC struct {
	storage Storage
	clock   cache.Clock
}

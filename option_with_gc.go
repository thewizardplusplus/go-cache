package cache

import (
	"github.com/thewizardplusplus/go-cache/gc"
	"github.com/thewizardplusplus/go-cache/models"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

// GCFactory ...
type GCFactory func(storage hashmap.Storage, clock models.Clock) gc.GC

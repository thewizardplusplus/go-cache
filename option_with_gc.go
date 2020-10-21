package cache

import (
	"context"
	"time"

	"github.com/thewizardplusplus/go-cache/gc"
	"github.com/thewizardplusplus/go-cache/models"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

// GCFactory ...
type GCFactory func(storage hashmap.Storage, clock models.Clock) gc.GC

// ConfigWithGC ...
type ConfigWithGC struct {
	ctx       context.Context
	storage   hashmap.Storage
	clock     models.Clock
	gcFactory GCFactory
	gcPeriod  time.Duration
}

// OptionWithGC ...
type OptionWithGC func(config *ConfigWithGC)

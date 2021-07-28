package cache

import (
	"time"

	"github.com/thewizardplusplus/go-cache/gc"
	"github.com/thewizardplusplus/go-cache/models"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

// GCFactory ...
type GCFactory func(storage hashmap.Storage, clock models.Clock) gc.GC

// ConfigWithGC ...
type ConfigWithGC struct {
	storage   hashmap.Storage
	clock     models.Clock
	gcFactory GCFactory
	gcPeriod  time.Duration
}

// OptionWithGC ...
type OptionWithGC func(config *ConfigWithGC)

// WithGCAndStorage ...
//
// It should be safe for concurrent access.
//
// Default: an instance of the hashmap.ConcurrentHashMap structure
// with default options.
//
func WithGCAndStorage(storage hashmap.Storage) OptionWithGC {
	return func(config *ConfigWithGC) {
		config.storage = storage
	}
}

// WithGCAndClock ...
//
// Default: the time.Now() function.
//
func WithGCAndClock(clock models.Clock) OptionWithGC {
	return func(config *ConfigWithGC) {
		config.clock = clock
	}
}

// WithGCAndGCFactory ...
//
// Default: a factory that produces an instance of the gc.PartialGC structure
// with default options.
//
func WithGCAndGCFactory(gcFactory GCFactory) OptionWithGC {
	return func(config *ConfigWithGC) {
		config.gcFactory = gcFactory
	}
}

// WithGCAndGCPeriod ...
//
// Default: 100 ms.
//
func WithGCAndGCPeriod(gcPeriod time.Duration) OptionWithGC {
	return func(config *ConfigWithGC) {
		config.gcPeriod = gcPeriod
	}
}

func newConfigWithGC(options []OptionWithGC) ConfigWithGC {
	// default config
	config := ConfigWithGC{
		storage: hashmap.NewConcurrentHashMap(),
		clock:   time.Now,
		gcFactory: func(storage hashmap.Storage, clock models.Clock) gc.GC {
			return gc.NewPartialGC(storage, gc.PartialGCWithClock(clock))
		},
		gcPeriod: 100 * time.Millisecond,
	}
	for _, option := range options {
		option(&config)
	}

	return config
}

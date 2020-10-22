package cache

import (
	"context"

	"github.com/thewizardplusplus/go-cache/gc"
	"github.com/thewizardplusplus/go-cache/models"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

//go:generate mockery -name=Context -inpkg -case=underscore -testonly

// Context ...
//
// It's used only for mock generating.
//
type Context interface {
	context.Context
}

//go:generate mockery -name=Key -inpkg -case=underscore -testonly

// Key ...
//
// It's used only for mock generating.
//
type Key interface {
	hashmap.Key
}

//go:generate mockery -name=Handler -inpkg -case=underscore -testonly

// Handler ...
//
// It's used only for mock generating.
//
type Handler interface {
	Handle(key hashmap.Key, value interface{}) bool
}

//go:generate mockery -name=Storage -inpkg -case=underscore -testonly

// Storage ...
//
// It's used only for mock generating.
//
type Storage interface {
	hashmap.Storage
}

//go:generate mockery -name=GC -inpkg -case=underscore -testonly

// GC ...
//
// It's used only for mock generating.
//
type GC interface {
	gc.GC
}

//go:generate mockery -name=GCFactoryHandler -inpkg -case=underscore -testonly

// GCFactoryHandler ...
//
// It's used only for mock generating.
//
type GCFactoryHandler interface {
	NewGC(storage hashmap.Storage, clock models.Clock) gc.GC
}

package hashmaputils

import (
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

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
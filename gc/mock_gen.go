package gc

import (
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

//go:generate mockery -name=Key -inpkg -case=underscore -testonly

// Key ...
//
// It's used only for mock generating.
type Key interface {
	hashmap.Key
}

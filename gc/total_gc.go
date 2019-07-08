package gc

import (
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

// Storage ...
type Storage interface {
	Iterate(handler hashmap.Handler) bool
	Delete(key hashmap.Key)
}

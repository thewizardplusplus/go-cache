package cache

import (
	"time"

	hashmap "github.com/thewizardplusplus/go-hashmap"
)

//go:generate mockery -name=Storage -inpkg -case=underscore -testonly

// Storage ...
type Storage interface {
	Get(key hashmap.Key) (data interface{}, ok bool)
	Set(key hashmap.Key, data interface{})
	Delete(key hashmap.Key)
}

// Clock ...
type Clock func() time.Time

// Cache ...
type Cache struct {
	storage Storage
	clock   Clock
}

// NewCache ...
func NewCache(storage Storage, clock Clock) Cache {
	return Cache{storage, clock}
}

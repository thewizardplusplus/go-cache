package cache

import (
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

// Storage ...
type Storage interface {
	Get(key hashmap.Key) (data interface{}, ok bool)
	Set(key hashmap.Key, data interface{})
	Delete(key hashmap.Key)
}

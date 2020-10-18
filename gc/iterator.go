package gc

import (
	cache "github.com/thewizardplusplus/go-cache"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

type iterator struct {
	counter

	storage hashmap.Storage
	clock   cache.Clock
}

func newIterator(storage hashmap.Storage, clock cache.Clock) *iterator {
	return &iterator{storage: storage, clock: clock}
}

func (iterator *iterator) handleIteration(
	key hashmap.Key,
	value interface{},
) bool {
	if value.(cache.Value).IsExpired(iterator.clock) {
		iterator.storage.Delete(key)
		iterator.expired++
	}

	iterator.iterated++
	return !iterator.stopIterate()
}

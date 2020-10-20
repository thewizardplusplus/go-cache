package gc

import (
	"github.com/thewizardplusplus/go-cache/models"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

type iterator struct {
	counter

	storage hashmap.Storage
	clock   models.Clock
}

func newIterator(
	storage hashmap.Storage,
	clock models.Clock,
	maxIteratedCount int,
	minExpiredPercent float64,
) *iterator {
	return &iterator{
		counter: newCounter(maxIteratedCount, minExpiredPercent),

		storage: storage,
		clock:   clock,
	}
}

func (iterator *iterator) handleIteration(
	key hashmap.Key,
	value interface{},
) bool {
	if value.(models.Value).IsExpired(iterator.clock) {
		iterator.storage.Delete(key)
		iterator.expiredCount++
	}

	iterator.iteratedCount++
	return !iterator.stopIterate()
}

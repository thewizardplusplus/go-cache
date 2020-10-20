package gc

type counter struct {
	maxIteratedCount  int
	minExpiredPercent float64

	iterated int
	expired  int
}

const (
	maxIteratedCount  = 20
	minExpiredPercent = 0.25
)

func (counter counter) stopIterate() bool {
	return counter.iterated >= counter.maxIteratedCount
}

func (counter counter) stopClean() bool {
	expiredPercent := float64(counter.expired) / float64(counter.iterated)
	return counter.iterated == 0 || expiredPercent < counter.minExpiredPercent
}

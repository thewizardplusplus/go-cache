package gc

type counter struct {
	iterated int
	expired  int
}

const (
	maxIteratedCount  = 20
	minExpiredPercent = 0.25
)

func (counter counter) stopIterate() bool {
	return counter.iterated >= maxIteratedCount
}

func (counter counter) stopClean() bool {
	expiredPercent := float64(counter.expired) / float64(counter.iterated)
	return counter.iterated == 0 || expiredPercent < minExpiredPercent
}

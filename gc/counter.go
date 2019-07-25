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

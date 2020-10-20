package gc

type counter struct {
	maxIteratedCount  int
	minExpiredPercent float64

	iteratedCount int
	expiredCount  int
}

const (
	maxIteratedCount  = 20
	minExpiredPercent = 0.25
)

func newCounter(maxIteratedCount int, minExpiredPercent float64) counter {
	return counter{
		maxIteratedCount:  maxIteratedCount,
		minExpiredPercent: minExpiredPercent,
	}
}

func (counter counter) stopIterate() bool {
	return counter.iteratedCount >= counter.maxIteratedCount
}

func (counter counter) stopClean() bool {
	expiredPercent :=
		float64(counter.expiredCount) / float64(counter.iteratedCount)
	return counter.iteratedCount == 0 || expiredPercent < counter.minExpiredPercent
}

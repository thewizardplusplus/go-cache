package gc

type counter struct {
	iterated int
	expired  int
}

const (
	maxIteratedCount  = 20
	minExpiredPercent = 0.25
)

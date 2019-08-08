package gc

import (
	"context"
	"time"
)

//go:generate mockery -name=GC -inpkg -case=underscore -testonly

// GC ...
type GC interface {
	Clean(ctx context.Context)
}

// Run ...
func Run(ctx context.Context, gc GC, period time.Duration) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			gc.Clean(ctx)
		case <-ctx.Done():
			return
		}
	}
}

package gc

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	cache "github.com/thewizardplusplus/go-cache"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

func BenchmarkCacheGetting_withPartialGC(benchmark *testing.B) {
	for _, expiredPercent := range []float32{0.01, 0.2, 0.3, 0.99} {
		for _, data := range []struct {
			name      string
			prepare   func(cache cache.Cache)
			benchmark func(cache cache.Cache)
		}{
			{
				name: "Get",
				prepare: func(cache cache.Cache) {
					for i := 0; i < sizeForBench; i++ {
						setItem(cache, i, expiredPercent)
					}
				},
				benchmark: func(cache cache.Cache) {
					cache.Get(IntKey(rand.Intn(sizeForBench))) // nolint: errcheck
				},
			},
			{
				name: "GetWithGC",
				prepare: func(cache cache.Cache) {
					for i := 0; i < sizeForBench; i++ {
						setItem(cache, i, expiredPercent)
					}
				},
				benchmark: func(cache cache.Cache) {
					cache.GetWithGC(IntKey(rand.Intn(sizeForBench))) // nolint: errcheck
				},
			},
		} {
			name := fmt.Sprintf("%s/%.2f", data.name, expiredPercent)
			benchmark.Run(name, func(benchmark *testing.B) {
				storage := hashmap.NewConcurrentHashMap()
				cache := cache.NewCache(storage, time.Now)
				data.prepare(cache)

				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				gc := NewPartialGC(storage, time.Now)
				go Run(ctx, gc, periodForBench)

				// add concurrent load
				go func() {
					ticker := time.NewTicker(periodForBench)
					defer ticker.Stop()

					for {
						select {
						case <-ticker.C:
							setItem(cache, rand.Intn(sizeForBench), expiredPercent)
						case <-ctx.Done():
							return
						}
					}
				}()

				benchmark.ResetTimer()

				for i := 0; i < benchmark.N; i++ {
					data.benchmark(cache)
				}
			})
		}
	}
}

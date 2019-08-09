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
	for _, data := range []struct {
		name      string
		prepare   func(cache cache.Cache, storageSize int, expiredPercent float32)
		benchmark func(cache cache.Cache, storageSize int)
	}{
		{
			name: "Get",
			prepare: func(cache cache.Cache, storageSize int, expiredPercent float32) {
				for i := 0; i < storageSize; i++ {
					setItem(cache, i, expiredPercent)
				}
			},
			benchmark: func(cache cache.Cache, storageSize int) {
				cache.Get(IntKey(rand.Intn(storageSize))) // nolint: errcheck
			},
		},
		{
			name: "GetWithGC",
			prepare: func(cache cache.Cache, storageSize int, expiredPercent float32) {
				for i := 0; i < storageSize; i++ {
					setItem(cache, i, expiredPercent)
				}
			},
			benchmark: func(cache cache.Cache, storageSize int) {
				cache.GetWithGC(IntKey(rand.Intn(storageSize))) // nolint: errcheck
			},
		},
	} {
		for _, storageSize := range []int{1e2, 1e4, 1e6} {
			for _, expiredPercent := range []float32{0.01, 0.2, 0.3, 0.99} {
				name := fmt.Sprintf("%s/%d/%.2f", data.name, storageSize, expiredPercent)
				benchmark.Run(name, func(benchmark *testing.B) {
					storage := hashmap.NewConcurrentHashMap()
					cache := cache.NewCache(storage, time.Now)
					data.prepare(cache, storageSize, expiredPercent)

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
								setItem(cache, rand.Intn(storageSize), expiredPercent)
							case <-ctx.Done():
								return
							}
						}
					}()

					benchmark.ResetTimer()

					for i := 0; i < benchmark.N; i++ {
						data.benchmark(cache, storageSize)
					}
				})
			}
		}
	}
}

package gc

import (
	"context"
	"math/rand"
	"testing"
	"time"

	cache "github.com/thewizardplusplus/go-cache"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

func BenchmarkCacheGetting_withPartialGC(benchmark *testing.B) {
	for _, data := range []struct {
		name      string
		prepare   func(cache cache.Cache)
		benchmark func(cache cache.Cache)
	}{
		{
			name: "Get",
			prepare: func(cache cache.Cache) {
				for i := 0; i < sizeForBench; i++ {
					setItem(cache, i)
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
					setItem(cache, i)
				}
			},
			benchmark: func(cache cache.Cache) {
				cache.GetWithGC(IntKey(rand.Intn(sizeForBench))) // nolint: errcheck
			},
		},
	} {
		benchmark.Run(data.name, func(benchmark *testing.B) {
			storage := hashmap.NewConcurrentHashMap()
			cache := cache.NewCache(storage, time.Now)
			data.prepare(cache)

			gc := NewPartialGC(storage, time.Now)
			go Run(context.Background(), gc, time.Nanosecond)

			// add concurrent load
			go func() {
				for {
					setItem(cache, rand.Intn(sizeForBench))
				}
			}()

			benchmark.ResetTimer()

			for i := 0; i < benchmark.N; i++ {
				data.benchmark(cache)
			}
		})
	}
}

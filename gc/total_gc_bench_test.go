package gc

import (
	"context"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math/rand"
	"testing"
	"time"

	cache "github.com/thewizardplusplus/go-cache"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

type IntKey int

func (key IntKey) Hash() int {
	hash := fnv.New32()
	binary.Write(hash, binary.LittleEndian, int32(key)) // nolint: errcheck

	return int(hash.Sum32())
}

func (key IntKey) Equals(other interface{}) bool {
	return key == other.(IntKey)
}

const (
	periodForBench = time.Nanosecond
)

func BenchmarkCacheGetting_withTotalGC(benchmark *testing.B) {
	for _, storageSize := range []int{1e2, 1e4, 1e6} {
		for _, expiredPercent := range []float32{0.01, 0.2, 0.3, 0.99} {
			for _, data := range []struct {
				name      string
				prepare   func(cache cache.Cache)
				benchmark func(cache cache.Cache)
			}{
				{
					name: "Get",
					prepare: func(cache cache.Cache) {
						for i := 0; i < storageSize; i++ {
							setItem(cache, i, expiredPercent)
						}
					},
					benchmark: func(cache cache.Cache) {
						cache.Get(IntKey(rand.Intn(storageSize))) // nolint: errcheck
					},
				},
				{
					name: "GetWithGC",
					prepare: func(cache cache.Cache) {
						for i := 0; i < storageSize; i++ {
							setItem(cache, i, expiredPercent)
						}
					},
					benchmark: func(cache cache.Cache) {
						cache.GetWithGC(IntKey(rand.Intn(storageSize))) // nolint: errcheck
					},
				},
			} {
				name := fmt.Sprintf("%s/%d/%.2f", data.name, storageSize, expiredPercent)
				benchmark.Run(name, func(benchmark *testing.B) {
					storage := hashmap.NewConcurrentHashMap()
					cache := cache.NewCache(storage, time.Now)
					data.prepare(cache)

					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()

					gc := NewTotalGC(storage, time.Now)
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
						data.benchmark(cache)
					}
				})
			}
		}
	}
}

func setItem(cache cache.Cache, key int, expiredPercent float32) {
	var ttl time.Duration
	// part of items will be already expired
	if rand.Float32() < expiredPercent {
		ttl = -time.Second
	}

	cache.Set(IntKey(key), key, ttl)
}

package gc_test

import (
	"context"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math/rand"
	"testing"
	"time"

	cache "github.com/thewizardplusplus/go-cache"
	"github.com/thewizardplusplus/go-cache/gc"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

const (
	periodForBench = time.Nanosecond
)

type IntKey int

func (key IntKey) Hash() int {
	hash := fnv.New32()
	binary.Write(hash, binary.LittleEndian, int32(key)) // nolint: errcheck

	return int(hash.Sum32())
}

func (key IntKey) Equals(other hashmap.Key) bool {
	return key == other.(IntKey)
}

func BenchmarkCacheGetting_withTotalGC(benchmark *testing.B) {
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
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()

					storage := hashmap.NewConcurrentHashMap()
					gcInstance := gc.NewTotalGC(storage)
					go gc.Run(ctx, gcInstance, periodForBench)

					cache := cache.NewCache(cache.WithStorage(storage))
					data.prepare(cache, storageSize, expiredPercent)

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

func setItem(cache cache.Cache, key int, expiredPercent float32) {
	var ttl time.Duration
	// part of items will be already expired
	if rand.Float32() < expiredPercent {
		ttl = -time.Second
	}

	cache.Set(IntKey(key), key, ttl)
}

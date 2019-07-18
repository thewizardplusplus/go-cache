package gc

import (
	"context"
	"encoding/binary"
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
	sizeForBench   = 1000
	periodForBench = time.Nanosecond
)

func BenchmarkCacheGetting_withTotalGC(benchmark *testing.B) {
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

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			gc := NewTotalGC(storage, time.Now)
			go Run(ctx, gc, periodForBench)

			// add concurrent load
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						setItem(cache, rand.Intn(sizeForBench))
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

func setItem(cache cache.Cache, key int) {
	var ttl time.Duration
	// half of items will be already expired
	if rand.Float32() < 0.5 {
		ttl = -time.Second
	}

	cache.Set(IntKey(key), key, ttl)
}

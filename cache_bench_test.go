package cache

import (
	"context"
	"encoding/binary"
	"hash/fnv"
	"math/rand"
	"testing"
	"time"

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
	expiredPercent = 0.5
)

func BenchmarkCacheGetting(benchmark *testing.B) {
	for _, data := range []struct {
		name      string
		prepare   func(cache Cache)
		benchmark func(cache Cache)
	}{
		{
			name: "Get",
			prepare: func(cache Cache) {
				for i := 0; i < sizeForBench; i++ {
					setItem(cache, i)
				}
			},
			benchmark: func(cache Cache) {
				cache.Get(IntKey(rand.Intn(sizeForBench))) // nolint: errcheck
			},
		},
		{
			name: "GetWithGC",
			prepare: func(cache Cache) {
				for i := 0; i < sizeForBench; i++ {
					setItem(cache, i)
				}
			},
			benchmark: func(cache Cache) {
				cache.GetWithGC(IntKey(rand.Intn(sizeForBench))) // nolint: errcheck
			},
		},
	} {
		benchmark.Run(data.name, func(benchmark *testing.B) {
			storage := hashmap.NewConcurrentHashMap()
			cache := NewCache(storage, time.Now)
			data.prepare(cache)

			// add concurrent load
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func() {
				ticker := time.NewTicker(periodForBench)
				defer ticker.Stop()

				for {
					select {
					case <-ticker.C:
						setItem(cache, rand.Intn(sizeForBench))
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

func setItem(cache Cache, key int) {
	var ttl time.Duration
	// part of items will be already expired
	if rand.Float32() < expiredPercent {
		ttl = -time.Second
	}

	cache.Set(IntKey(key), key, ttl)
}

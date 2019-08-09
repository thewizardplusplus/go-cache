package cache

import (
	"context"
	"encoding/binary"
	"fmt"
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

func BenchmarkCacheGetting(benchmark *testing.B) {
	for _, storageSize := range []int{1e2, 1e4, 1e6} {
		for _, data := range []struct {
			name      string
			prepare   func(cache Cache)
			benchmark func(cache Cache)
		}{
			{
				name: "Get",
				prepare: func(cache Cache) {
					for i := 0; i < storageSize; i++ {
						setItem(cache, i)
					}
				},
				benchmark: func(cache Cache) {
					cache.Get(IntKey(rand.Intn(storageSize))) // nolint: errcheck
				},
			},
			{
				name: "GetWithGC",
				prepare: func(cache Cache) {
					for i := 0; i < storageSize; i++ {
						setItem(cache, i)
					}
				},
				benchmark: func(cache Cache) {
					cache.GetWithGC(IntKey(rand.Intn(storageSize))) // nolint: errcheck
				},
			},
		} {
			name := fmt.Sprintf("%s/%d", data.name, storageSize)
			benchmark.Run(name, func(benchmark *testing.B) {
				storage := hashmap.NewConcurrentHashMap()
				cache := NewCache(storage, time.Now)
				data.prepare(cache)

				// add concurrent load
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				go func() {
					ticker := time.NewTicker(time.Nanosecond)
					defer ticker.Stop()

					for {
						select {
						case <-ticker.C:
							setItem(cache, rand.Intn(storageSize))
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

func setItem(cache Cache, key int) {
	var ttl time.Duration
	// half of items will be already expired
	if rand.Float32() < 0.5 {
		ttl = -time.Second
	}

	cache.Set(IntKey(key), key, ttl)
}

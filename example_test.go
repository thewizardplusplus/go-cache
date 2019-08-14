package cache_test

import (
	"context"
	"fmt"
	"hash/fnv"
	"io"
	"time"

	cache "github.com/thewizardplusplus/go-cache"
	"github.com/thewizardplusplus/go-cache/gc"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

type StringKey string

func (key StringKey) Hash() int {
	hash := fnv.New32()
	io.WriteString(hash, string(key)) // nolint: errcheck

	return int(hash.Sum32())
}

func (key StringKey) Equals(other interface{}) bool {
	return key == other.(StringKey)
}

const (
	gcPeriod     = time.Millisecond
	exampleDelay = gcPeriod * 100
)

func Example() {
	storage := hashmap.NewConcurrentHashMap()
	cleaner := gc.NewPartialGC(storage, time.Now)
	go gc.Run(context.Background(), cleaner, gcPeriod)

	timeZones := cache.NewCache(storage, time.Now)
	timeZones.Set(StringKey("EST"), -5*60*60, exampleDelay/2)
	timeZones.Set(StringKey("CST"), -6*60*60, exampleDelay/2)
	timeZones.Set(StringKey("MST"), -7*60*60, exampleDelay/2)

	estOffset, err := timeZones.Get(StringKey("EST"))
	fmt.Println(estOffset, err)

	time.Sleep(exampleDelay)

	estOffset, err = timeZones.Get(StringKey("EST"))
	fmt.Println(estOffset, err)

	// Output:
	// -18000 <nil>
	// <nil> key missed
}

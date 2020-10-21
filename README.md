# go-cache

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-cache?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-cache)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-cache)](https://goreportcard.com/report/github.com/thewizardplusplus/go-cache)
[![Build Status](https://travis-ci.org/thewizardplusplus/go-cache.svg?branch=master)](https://travis-ci.org/thewizardplusplus/go-cache)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-cache/branch/master/graph/badge.svg)](https://codecov.io/gh/thewizardplusplus/go-cache)

## Installation

```
$ go get github.com/thewizardplusplus/go-cache
```

## Example

`cache.NewCache()`:

```go
package main

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

func (key StringKey) Equals(other hashmap.Key) bool {
	return key == other.(StringKey)
}

const (
	gcPeriod     = time.Millisecond
	exampleDelay = gcPeriod * 100
)

func main() {
	storage := hashmap.NewConcurrentHashMap()
	gcInstance := gc.NewPartialGC(storage)
	go gc.Run(context.Background(), gcInstance, gcPeriod)

	timeZones := cache.NewCache(cache.WithStorage(storage))
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
```

`cache.NewCacheWithGC()`:

```go
package main

import (
	"fmt"
	"hash/fnv"
	"io"
	"time"

	cache "github.com/thewizardplusplus/go-cache"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

type StringKey string

func (key StringKey) Hash() int {
	hash := fnv.New32()
	io.WriteString(hash, string(key)) // nolint: errcheck

	return int(hash.Sum32())
}

func (key StringKey) Equals(other hashmap.Key) bool {
	return key == other.(StringKey)
}

const (
	gcPeriod     = time.Millisecond
	exampleDelay = gcPeriod * 100
)

func main() {
	timeZones := cache.NewCacheWithGC(cache.WithGCAndGCPeriod(gcPeriod))
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
```

## Benchmarks

```
BenchmarkCacheGetting/Get/100-8                   	30000000	       454 ns/op	      41 B/op	       7 allocs/op
BenchmarkCacheGetting/Get/10000-8                 	30000000	       638 ns/op	      42 B/op	       7 allocs/op
BenchmarkCacheGetting/Get/1000000-8               	10000000	      1583 ns/op	      46 B/op	       7 allocs/op
BenchmarkCacheGetting/GetWithGC/100-8             	30000000	       403 ns/op	      41 B/op	       7 allocs/op
BenchmarkCacheGetting/GetWithGC/10000-8           	30000000	       515 ns/op	      42 B/op	       7 allocs/op
BenchmarkCacheGetting/GetWithGC/1000000-8         	20000000	       934 ns/op	      44 B/op	       7 allocs/op
```

With the total GC:

```
BenchmarkCacheGetting_withTotalGC/Get/100/0.01-8                     	20000000	       832 ns/op	      94 B/op	       7 allocs/op
BenchmarkCacheGetting_withTotalGC/Get/100/0.20-8                     	20000000	       801 ns/op	      89 B/op	       7 allocs/op
BenchmarkCacheGetting_withTotalGC/Get/100/0.30-8                     	20000000	       859 ns/op	      90 B/op	       7 allocs/op
BenchmarkCacheGetting_withTotalGC/Get/100/0.99-8                     	20000000	       798 ns/op	      88 B/op	      10 allocs/op
BenchmarkCacheGetting_withTotalGC/Get/10000/0.01-8                   	10000000	      1112 ns/op	      94 B/op	       7 allocs/op
BenchmarkCacheGetting_withTotalGC/Get/10000/0.20-8                   	10000000	      1158 ns/op	     107 B/op	       7 allocs/op
BenchmarkCacheGetting_withTotalGC/Get/10000/0.30-8                   	20000000	      1016 ns/op	     103 B/op	       7 allocs/op
BenchmarkCacheGetting_withTotalGC/Get/10000/0.99-8                   	30000000	       466 ns/op	      54 B/op	       8 allocs/op
BenchmarkCacheGetting_withTotalGC/Get/1000000/0.01-8                 	10000000	      1959 ns/op	     114 B/op	       7 allocs/op
BenchmarkCacheGetting_withTotalGC/Get/1000000/0.20-8                 	10000000	      1685 ns/op	      92 B/op	       7 allocs/op
BenchmarkCacheGetting_withTotalGC/Get/1000000/0.30-8                 	10000000	      1365 ns/op	      76 B/op	       8 allocs/op
BenchmarkCacheGetting_withTotalGC/Get/1000000/0.99-8                 	20000000	       664 ns/op	      54 B/op	       8 allocs/op
BenchmarkCacheGetting_withTotalGC/GetWithGC/100/0.01-8               	20000000	       793 ns/op	      95 B/op	       7 allocs/op
BenchmarkCacheGetting_withTotalGC/GetWithGC/100/0.20-8               	20000000	       794 ns/op	      90 B/op	       7 allocs/op
BenchmarkCacheGetting_withTotalGC/GetWithGC/100/0.30-8               	20000000	       833 ns/op	      91 B/op	       7 allocs/op
BenchmarkCacheGetting_withTotalGC/GetWithGC/100/0.99-8               	20000000	       660 ns/op	      76 B/op	       8 allocs/op
BenchmarkCacheGetting_withTotalGC/GetWithGC/10000/0.01-8             	20000000	       981 ns/op	      92 B/op	       7 allocs/op
BenchmarkCacheGetting_withTotalGC/GetWithGC/10000/0.20-8             	20000000	       947 ns/op	     104 B/op	       7 allocs/op
BenchmarkCacheGetting_withTotalGC/GetWithGC/10000/0.30-8             	20000000	       937 ns/op	     105 B/op	       7 allocs/op
BenchmarkCacheGetting_withTotalGC/GetWithGC/10000/0.99-8             	30000000	       530 ns/op	      56 B/op	       8 allocs/op
BenchmarkCacheGetting_withTotalGC/GetWithGC/1000000/0.01-8           	10000000	      2000 ns/op	     114 B/op	       7 allocs/op
BenchmarkCacheGetting_withTotalGC/GetWithGC/1000000/0.20-8           	10000000	      1859 ns/op	      95 B/op	       8 allocs/op
BenchmarkCacheGetting_withTotalGC/GetWithGC/1000000/0.30-8           	10000000	      1428 ns/op	      80 B/op	       8 allocs/op
BenchmarkCacheGetting_withTotalGC/GetWithGC/1000000/0.99-8           	20000000	       686 ns/op	      59 B/op	       8 allocs/op
```

With the partial GC:

```
BenchmarkCacheGetting_withPartialGC/Get/100/0.01-8                   	30000000	       449 ns/op	      46 B/op	       7 allocs/op
BenchmarkCacheGetting_withPartialGC/Get/100/0.20-8                   	30000000	       433 ns/op	      46 B/op	       7 allocs/op
BenchmarkCacheGetting_withPartialGC/Get/100/0.30-8                   	30000000	       446 ns/op	      46 B/op	       7 allocs/op
BenchmarkCacheGetting_withPartialGC/Get/100/0.99-8                   	20000000	       895 ns/op	     113 B/op	      11 allocs/op
BenchmarkCacheGetting_withPartialGC/Get/10000/0.01-8                 	10000000	      1564 ns/op	     237 B/op	       7 allocs/op
BenchmarkCacheGetting_withPartialGC/Get/10000/0.20-8                 	10000000	      1807 ns/op	     256 B/op	       7 allocs/op
BenchmarkCacheGetting_withPartialGC/Get/10000/0.30-8                 	10000000	      1473 ns/op	     239 B/op	       7 allocs/op
BenchmarkCacheGetting_withPartialGC/Get/10000/0.99-8                 	20000000	       839 ns/op	     130 B/op	       8 allocs/op
BenchmarkCacheGetting_withPartialGC/Get/1000000/0.01-8               	 1000000	     12194 ns/op	    2462 B/op	       7 allocs/op
BenchmarkCacheGetting_withPartialGC/Get/1000000/0.20-8               	 1000000	     12829 ns/op	    2540 B/op	       7 allocs/op
BenchmarkCacheGetting_withPartialGC/Get/1000000/0.30-8               	 1000000	     12532 ns/op	    2413 B/op	       7 allocs/op
BenchmarkCacheGetting_withPartialGC/Get/1000000/0.99-8               	 1000000	     11354 ns/op	    2164 B/op	       7 allocs/op
BenchmarkCacheGetting_withPartialGC/GetWithGC/100/0.01-8             	30000000	       460 ns/op	      46 B/op	       7 allocs/op
BenchmarkCacheGetting_withPartialGC/GetWithGC/100/0.20-8             	30000000	       446 ns/op	      46 B/op	       7 allocs/op
BenchmarkCacheGetting_withPartialGC/GetWithGC/100/0.30-8             	30000000	       473 ns/op	      46 B/op	       7 allocs/op
BenchmarkCacheGetting_withPartialGC/GetWithGC/100/0.99-8             	20000000	       978 ns/op	     120 B/op	      11 allocs/op
BenchmarkCacheGetting_withPartialGC/GetWithGC/10000/0.01-8           	10000000	      1594 ns/op	     243 B/op	       7 allocs/op
BenchmarkCacheGetting_withPartialGC/GetWithGC/10000/0.20-8           	10000000	      1554 ns/op	     244 B/op	       7 allocs/op
BenchmarkCacheGetting_withPartialGC/GetWithGC/10000/0.30-8           	10000000	      1548 ns/op	     240 B/op	       7 allocs/op
BenchmarkCacheGetting_withPartialGC/GetWithGC/10000/0.99-8           	20000000	       886 ns/op	     134 B/op	       8 allocs/op
BenchmarkCacheGetting_withPartialGC/GetWithGC/1000000/0.01-8         	 1000000	     15485 ns/op	    3050 B/op	       7 allocs/op
BenchmarkCacheGetting_withPartialGC/GetWithGC/1000000/0.20-8         	  300000	     37861 ns/op	    7055 B/op	       9 allocs/op
BenchmarkCacheGetting_withPartialGC/GetWithGC/1000000/0.30-8         	  300000	     49486 ns/op	    8154 B/op	       9 allocs/op
BenchmarkCacheGetting_withPartialGC/GetWithGC/1000000/0.99-8         	  200000	     55579 ns/op	    8959 B/op	      13 allocs/op
```

## License

The MIT License (MIT)

Copyright &copy; 2019-2020 thewizardplusplus

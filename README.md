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

```go
package main

import (
	"fmt"
	"hash/fnv"
	"io"
	"time"

	hashmap "github.com/thewizardplusplus/go-hashmap"
)

type StringKey string

func (key StringKey) Hash() int {
	hash := fnv.New32()
	io.WriteString(hash, string(key))

	return int(hash.Sum32())
}

func (key StringKey) Equals(other interface{}) bool {
	return key == other.(StringKey)
}

func main() {
	storage := hashmap.NewConcurrentHashMap()
	timeZones := NewCache(storage, time.Now)
	timeZones.Set(StringKey("EST"), -5*60*60, 100*time.Millisecond)
	timeZones.Set(StringKey("CST"), -6*60*60, 100*time.Millisecond)
	timeZones.Set(StringKey("MST"), -7*60*60, 100*time.Millisecond)

	estOffset, err := timeZones.Get(StringKey("EST"))
	fmt.Println(estOffset, err)

	time.Sleep(200 * time.Millisecond)

	estOffset, err = timeZones.Get(StringKey("EST"))
	fmt.Println(estOffset, err)

	// Output:
	// -18000 <nil>
	// <nil> key expired
}
```

## Benchmarks

```
BenchmarkCacheGetting/Get-8         	10000000	      1641 ns/op	     352 B/op	      36 allocs/op
BenchmarkCacheGetting/GetWithGC-8   	 5000000	      3291 ns/op	     575 B/op	      57 allocs/op
```

## License

The MIT License (MIT)

Copyright &copy; 2019 thewizardplusplus

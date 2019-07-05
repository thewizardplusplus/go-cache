package cache

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
	io.WriteString(hash, string(key)) // nolint: errcheck

	return int(hash.Sum32())
}

func (key StringKey) Equals(other interface{}) bool {
	return key == other.(StringKey)
}

func Example() {
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

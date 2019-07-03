package cache

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

func TestNewCache(test *testing.T) {
	storage := new(MockStorage)
	cache := NewCache(storage, time.Now)

	mock.AssertExpectationsForObjects(test, storage)
	assert.Equal(test, storage, cache.storage)
	assert.Equal(test, getPointer(time.Now), getPointer(cache.clock))
}

func TestCache_Set(test *testing.T) {
	type fields struct {
		storage Storage
		clock   Clock
	}
	type args struct {
		key  hashmap.Key
		data interface{}
		ttl  time.Duration
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			cache := Cache{data.fields.storage, data.fields.clock}
			cache.Set(data.args.key, data.args.data, data.args.ttl)

			mock.AssertExpectationsForObjects(test, data.fields.storage, data.args.key)
		})
	}
}

func getPointer(value interface{}) uintptr {
	return reflect.ValueOf(value).Pointer()
}

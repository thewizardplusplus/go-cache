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

func TestCache_Get(test *testing.T) {
	type fields struct {
		storage Storage
		clock   Clock
	}
	type args struct {
		key hashmap.Key
	}

	for _, data := range []struct {
		name     string
		fields   fields
		args     args
		wantData interface{}
		wantErr  assert.ErrorAssertionFunc
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			cache := Cache{data.fields.storage, data.fields.clock}
			gotData, gotErr := cache.Get(data.args.key)

			mock.AssertExpectationsForObjects(test, data.fields.storage, data.args.key)
			assert.Equal(test, data.wantData, gotData)
			data.wantErr(test, gotErr)
		})
	}
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
		{
			name: "without a TTL",
			fields: fields{
				storage: func() Storage {
					storage := new(MockStorage)
					storage.On("Set", new(MockKey), value{"data", time.Time{}})

					return storage
				}(),
				clock: clock,
			},
			args: args{
				key:  new(MockKey),
				data: "data",
				ttl:  0,
			},
		},
		{
			name: "with a TTL",
			fields: fields{
				storage: func() Storage {
					storage := new(MockStorage)
					storage.On("Set", new(MockKey), value{"data", clock().Add(time.Second)})

					return storage
				}(),
				clock: clock,
			},
			args: args{
				key:  new(MockKey),
				data: "data",
				ttl:  time.Second,
			},
		},
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

func clock() time.Time {
	return time.Date(
		2006, time.January, 2, // year, month, day
		15, 4, 5, // hour, minute, second
		0,        // nanosecond
		time.UTC, // location
	)
}

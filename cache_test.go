package cache

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

type MockKeyWithID struct {
	MockKey

	ID int
}

func NewMockKeyWithID(id int) *MockKeyWithID {
	return &MockKeyWithID{ID: id}
}

func TestValue_IsExpired(test *testing.T) {
	type fields struct {
		Data           interface{}
		ExpirationTime time.Time
	}
	type args struct {
		clock Clock
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
		want   assert.BoolAssertionFunc
	}{
		{
			name: "zero expiration time",
			fields: fields{
				Data:           "data",
				ExpirationTime: time.Time{},
			},
			args: args{
				clock: clock,
			},
			want: assert.False,
		},
		{
			name: "expiration time less than the current one",
			fields: fields{
				Data:           "data",
				ExpirationTime: clock().Add(-time.Second),
			},
			args: args{
				clock: clock,
			},
			want: assert.True,
		},
		{
			name: "expiration time equal to the current one",
			fields: fields{
				Data:           "data",
				ExpirationTime: clock(),
			},
			args: args{
				clock: clock,
			},
			want: assert.False,
		},
		{
			name: "expiration time greater than the current one",
			fields: fields{
				Data:           "data",
				ExpirationTime: clock().Add(time.Second),
			},
			args: args{
				clock: clock,
			},
			want: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := Value{data.fields.Data, data.fields.ExpirationTime}
			got := value.IsExpired(data.args.clock)

			data.want(test, got)
		})
	}
}

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
		{
			name: "success with a zero expiration time",
			fields: fields{
				storage: func() Storage {
					storage := new(MockStorage)
					storage.
						On("Get", NewMockKeyWithID(23)).
						Return(Value{"data", time.Time{}}, true)

					return storage
				}(),
				clock: clock,
			},
			args: args{
				key: NewMockKeyWithID(23),
			},
			wantData: "data",
			wantErr:  assert.NoError,
		},
		{
			name: "success with an expiration time equal to current one",
			fields: fields{
				storage: func() Storage {
					storage := new(MockStorage)
					storage.
						On("Get", NewMockKeyWithID(23)).
						Return(Value{"data", clock()}, true)

					return storage
				}(),
				clock: clock,
			},
			args: args{
				key: NewMockKeyWithID(23),
			},
			wantData: "data",
			wantErr:  assert.NoError,
		},
		{
			name: "success with an expiration time greater than current one",
			fields: fields{
				storage: func() Storage {
					storage := new(MockStorage)
					storage.
						On("Get", NewMockKeyWithID(23)).
						Return(Value{"data", clock().Add(time.Second)}, true)

					return storage
				}(),
				clock: clock,
			},
			args: args{
				key: NewMockKeyWithID(23),
			},
			wantData: "data",
			wantErr:  assert.NoError,
		},
		{
			name: "error with a missed key",
			fields: fields{
				storage: func() Storage {
					storage := new(MockStorage)
					storage.On("Get", NewMockKeyWithID(23)).Return(nil, false)

					return storage
				}(),
				clock: clock,
			},
			args: args{
				key: NewMockKeyWithID(23),
			},
			wantData: nil,
			wantErr:  assert.Error,
		},
		{
			name: "error with an expired key",
			fields: fields{
				storage: func() Storage {
					storage := new(MockStorage)
					storage.
						On("Get", NewMockKeyWithID(23)).
						Return(Value{"data", clock().Add(-time.Second)}, true)

					return storage
				}(),
				clock: clock,
			},
			args: args{
				key: NewMockKeyWithID(23),
			},
			wantData: nil,
			wantErr:  assert.Error,
		},
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

func TestCache_GetWithGC(test *testing.T) {
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
		{
			name: "success",
			fields: fields{
				storage: func() Storage {
					storage := new(MockStorage)
					storage.
						On("Get", NewMockKeyWithID(23)).
						Return(Value{"data", clock().Add(time.Second)}, true)

					return storage
				}(),
				clock: clock,
			},
			args: args{
				key: NewMockKeyWithID(23),
			},
			wantData: "data",
			wantErr:  assert.NoError,
		},
		{
			name: "error with a missed key",
			fields: fields{
				storage: func() Storage {
					storage := new(MockStorage)
					storage.On("Get", NewMockKeyWithID(23)).Return(nil, false)

					return storage
				}(),
				clock: clock,
			},
			args: args{
				key: NewMockKeyWithID(23),
			},
			wantData: nil,
			wantErr:  assert.Error,
		},
		{
			name: "error with an expired key",
			fields: fields{
				storage: func() Storage {
					storage := new(MockStorage)
					storage.
						On("Get", NewMockKeyWithID(23)).
						Return(Value{"data", clock().Add(-time.Second)}, true)
					storage.On("Delete", NewMockKeyWithID(23))

					return storage
				}(),
				clock: clock,
			},
			args: args{
				key: NewMockKeyWithID(23),
			},
			wantData: nil,
			wantErr:  assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			cache := Cache{data.fields.storage, data.fields.clock}
			gotData, gotErr := cache.GetWithGC(data.args.key)

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
					storage.On("Set", NewMockKeyWithID(23), Value{"data", time.Time{}})

					return storage
				}(),
				clock: clock,
			},
			args: args{
				key:  NewMockKeyWithID(23),
				data: "data",
				ttl:  0,
			},
		},
		{
			name: "with a TTL",
			fields: fields{
				storage: func() Storage {
					storage := new(MockStorage)
					storage.
						On("Set", NewMockKeyWithID(23), Value{"data", clock().Add(time.Second)})

					return storage
				}(),
				clock: clock,
			},
			args: args{
				key:  NewMockKeyWithID(23),
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

func TestCache_Delete(test *testing.T) {
	storage := new(MockStorage)
	storage.On("Delete", NewMockKeyWithID(23))

	key := NewMockKeyWithID(23)

	cache := Cache{storage, clock}
	cache.Delete(key)

	mock.AssertExpectationsForObjects(test, storage, key)
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

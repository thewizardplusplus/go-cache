package cache

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-cache/models"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

type MockKeyWithID struct {
	MockKey

	ID int
}

func NewMockKeyWithID(id int) *MockKeyWithID {
	return &MockKeyWithID{ID: id}
}

func (key *MockKeyWithID) Copy() *MockKeyWithID {
	return NewMockKeyWithID(key.ID)
}

func TestNewCache(test *testing.T) {
	type args struct {
		options []Option
	}

	for _, data := range []struct {
		name          string
		args          args
		wantStorage   hashmap.Storage
		wantClockTime time.Time
	}{
		{
			name: "with default options",
			args: args{
				options: nil,
			},
			wantStorage:   hashmap.NewConcurrentHashMap(),
			wantClockTime: time.Now(),
		},
		{
			name: "with the set storage",
			args: args{
				options: []Option{WithStorage(new(MockStorage))},
			},
			wantStorage:   new(MockStorage),
			wantClockTime: time.Now(),
		},
		{
			name: "with the set clock",
			args: args{
				options: []Option{WithClock(clock)},
			},
			wantStorage:   hashmap.NewConcurrentHashMap(),
			wantClockTime: clock(),
		},
		{
			name: "with set options",
			args: args{
				options: []Option{WithStorage(new(MockStorage)), WithClock(clock)},
			},
			wantStorage:   new(MockStorage),
			wantClockTime: clock(),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := NewCache(data.args.options...)

			_, ok := got.storage.(interface {
				AssertExpectations(assert.TestingT) bool // nolint: staticcheck
			})
			if ok {
				mock.AssertExpectationsForObjects(test, got.storage)
			}
			assert.Equal(test, data.wantStorage, got.storage)

			// don't use the reflect.Value.Pointer() method for this check; see details:
			// * https://golang.org/pkg/reflect/#Value.Pointer
			// * https://stackoverflow.com/a/9644797
			require.NotNil(test, got.clock)
			assert.WithinDuration(test, data.wantClockTime, got.clock(), time.Hour)
		})
	}
}

func TestNewCacheWithGC(test *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	storage := new(MockStorage)

	gcInstance := new(MockGC)
	gcInstance.On("Clean", ctx).Return()

	gcFactoryHandler := new(MockGCFactoryHandler)
	gcFactoryHandler.
		On("NewGC", storage, mock.AnythingOfType("models.Clock")).
		Return(gcInstance)

	const gcPeriod = 100 * time.Millisecond
	cache := NewCacheWithGC(
		WithGCAndContext(ctx),
		WithGCAndStorage(storage),
		WithGCAndClock(clock),
		WithGCAndGCFactory(gcFactoryHandler.NewGC),
		WithGCAndGCPeriod(gcPeriod),
	)

	time.Sleep(gcPeriod * 2)
	cancel()

	mock.AssertExpectationsForObjects(test, storage, gcInstance, gcFactoryHandler)
	assert.Equal(test, storage, cache.storage)

	// don't use the reflect.Value.Pointer() method for this check; see details:
	// * https://golang.org/pkg/reflect/#Value.Pointer
	// * https://stackoverflow.com/a/9644797
	require.NotNil(test, cache.clock)
	assert.WithinDuration(test, clock(), cache.clock(), time.Hour)
}

func TestCache_Get(test *testing.T) {
	type fields struct {
		storage hashmap.Storage
		clock   models.Clock
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
				storage: func() hashmap.Storage {
					storage := new(MockStorage)
					storage.
						On("Get", NewMockKeyWithID(23)).
						Return(
							models.Value{Data: "data", ExpirationTime: clock().Add(time.Second)},
							true,
						)

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
				storage: func() hashmap.Storage {
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
				storage: func() hashmap.Storage {
					storage := new(MockStorage)
					storage.
						On("Get", NewMockKeyWithID(23)).
						Return(
							models.Value{Data: "data", ExpirationTime: clock().Add(-time.Second)},
							true,
						)

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
		storage hashmap.Storage
		clock   models.Clock
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
				storage: func() hashmap.Storage {
					storage := new(MockStorage)
					storage.
						On("Get", NewMockKeyWithID(23)).
						Return(
							models.Value{Data: "data", ExpirationTime: clock().Add(time.Second)},
							true,
						)

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
				storage: func() hashmap.Storage {
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
				storage: func() hashmap.Storage {
					storage := new(MockStorage)
					storage.
						On("Get", NewMockKeyWithID(23)).
						Return(
							models.Value{Data: "data", ExpirationTime: clock().Add(-time.Second)},
							true,
						)
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

func TestCache_Iterate(test *testing.T) {
	type bucket struct {
		key   hashmap.Key
		value interface{}
		ttl   time.Duration
	}
	type fields struct {
		buckets []bucket
		clock   models.Clock
	}

	for _, data := range []struct {
		name             string
		fields           fields
		interruptOnCount int
		wantBuckets      []bucket
		wantOk           assert.BoolAssertionFunc
	}{
		{
			name: "without buckets",
			fields: fields{
				buckets: nil,
				clock:   clock,
			},
			interruptOnCount: 10,
			wantBuckets:      nil,
			wantOk:           assert.True,
		},
		{
			name: "with few not expired buckets and without an interrupt",
			fields: fields{
				buckets: func() []bucket {
					keyOne := NewMockKeyWithID(12)
					keyOne.On("Hash").Return(12)

					keyTwo := NewMockKeyWithID(23)
					keyTwo.On("Hash").Return(23)

					keyThree := NewMockKeyWithID(42)
					keyThree.On("Hash").Return(42)

					return []bucket{
						{key: keyOne, value: "one"},
						{key: keyTwo, value: "two"},
						{key: keyThree, value: "three"},
					}
				}(),
				clock: clock,
			},
			interruptOnCount: 10,
			wantBuckets: []bucket{
				{key: NewMockKeyWithID(12), value: "one"},
				{key: NewMockKeyWithID(42), value: "three"},
				{key: NewMockKeyWithID(23), value: "two"},
			},
			wantOk: assert.True,
		},
		{
			name: "with few not expired buckets and with an interrupt",
			fields: fields{
				buckets: func() []bucket {
					keyOne := NewMockKeyWithID(12)
					keyOne.On("Hash").Return(12)

					keyTwo := NewMockKeyWithID(23)
					keyTwo.On("Hash").Return(23)

					keyThree := NewMockKeyWithID(42)
					keyThree.On("Hash").Return(42)

					return []bucket{
						{key: keyOne, value: "one"},
						{key: keyTwo, value: "two"},
						{key: keyThree, value: "three"},
					}
				}(),
				clock: clock,
			},
			interruptOnCount: 2,
			wantBuckets: []bucket{
				{key: NewMockKeyWithID(12), value: "one"},
				{key: NewMockKeyWithID(42), value: "three"},
			},
			wantOk: assert.False,
		},
		{
			name: "with few not expired and expired buckets and without an interrupt",
			fields: fields{
				buckets: func() []bucket {
					keyOne := NewMockKeyWithID(12)
					keyOne.On("Hash").Return(12)

					keyTwo := NewMockKeyWithID(23)
					keyTwo.On("Hash").Return(23)

					keyThree := NewMockKeyWithID(42)
					keyThree.On("Hash").Return(42)

					return []bucket{
						{key: keyOne, value: "one"},
						{key: keyTwo, value: "two", ttl: -time.Second},
						{key: keyThree, value: "three", ttl: -time.Second},
					}
				}(),
				clock: clock,
			},
			interruptOnCount: 10,
			wantBuckets: []bucket{
				{key: NewMockKeyWithID(12), value: "one"},
			},
			wantOk: assert.True,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			// reset the random generator to make tests deterministic
			rand.Seed(1)

			cache := Cache{
				storage: hashmap.NewConcurrentHashMap(),
				clock:   data.fields.clock,
			}
			for _, bucket := range data.fields.buckets {
				cache.Set(bucket.key, bucket.value, bucket.ttl)
			}

			var gotBuckets []bucket
			gotOk := cache.Iterate(func(key hashmap.Key, value interface{}) bool {
				gotBuckets = append(gotBuckets, bucket{
					key:   key.(*MockKeyWithID).Copy(),
					value: value,
				})

				// interrupt after a specified count of got buckets
				return len(gotBuckets) < data.interruptOnCount
			})

			for index, bucket := range data.fields.buckets {
				mock.AssertExpectationsForObjects(test, bucket.key)

				data.fields.buckets[index].key = bucket.key.(*MockKeyWithID).Copy()
			}
			assert.Equal(test, data.wantBuckets, gotBuckets)
			data.wantOk(test, gotOk)
		})
	}
}

func TestCache_IterateWithGC(test *testing.T) {
	type bucket struct {
		key   hashmap.Key
		value interface{}
		ttl   time.Duration
	}
	type fields struct {
		buckets []bucket
		clock   models.Clock
	}

	for _, data := range []struct {
		name             string
		fields           fields
		interruptOnCount int
		wantBuckets      []bucket
		wantOk           assert.BoolAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			// reset the random generator to make tests deterministic
			rand.Seed(1)

			cache := Cache{
				storage: hashmap.NewConcurrentHashMap(),
				clock:   data.fields.clock,
			}
			for _, bucket := range data.fields.buckets {
				cache.Set(bucket.key, bucket.value, bucket.ttl)
			}

			var gotBuckets []bucket
			gotOk := cache.IterateWithGC(func(key hashmap.Key, value interface{}) bool {
				gotBuckets = append(gotBuckets, bucket{
					key:   key.(*MockKeyWithID).Copy(),
					value: value,
				})

				// interrupt after a specified count of got buckets
				return len(gotBuckets) < data.interruptOnCount
			})

			for index, bucket := range data.fields.buckets {
				mock.AssertExpectationsForObjects(test, bucket.key)

				data.fields.buckets[index].key = bucket.key.(*MockKeyWithID).Copy()
			}
			assert.Equal(test, data.wantBuckets, gotBuckets)
			data.wantOk(test, gotOk)
		})
	}
}

func TestCache_Set(test *testing.T) {
	type fields struct {
		storage hashmap.Storage
		clock   models.Clock
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
				storage: func() hashmap.Storage {
					storage := new(MockStorage)
					storage.On("Set", NewMockKeyWithID(23), models.Value{
						Data:           "data",
						ExpirationTime: time.Time{},
					})

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
				storage: func() hashmap.Storage {
					storage := new(MockStorage)
					storage.
						On("Set", NewMockKeyWithID(23), models.Value{
							Data:           "data",
							ExpirationTime: clock().Add(time.Second),
						})

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

func clock() time.Time {
	return time.Date(
		2006, time.January, 2, // year, month, day
		15, 4, 5, // hour, minute, second
		0,        // nanosecond
		time.UTC, // location
	)
}

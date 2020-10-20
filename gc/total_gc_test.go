package gc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	cache "github.com/thewizardplusplus/go-cache"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

const (
	timedTestDelay = 100 * time.Millisecond
)

type MockKeyWithID struct {
	MockKey

	ID int
}

func NewMockKeyWithID(id int) *MockKeyWithID {
	return &MockKeyWithID{ID: id}
}

func TestNewTotalGC(test *testing.T) {
	type args struct {
		storage hashmap.Storage
		options []TotalGCOption
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
				storage: new(MockStorage),
				options: nil,
			},
			wantStorage:   new(MockStorage),
			wantClockTime: time.Now(),
		},
		{
			name: "with the set clock",
			args: args{
				storage: new(MockStorage),
				options: []TotalGCOption{TotalGCWithClock(clock)},
			},
			wantStorage:   new(MockStorage),
			wantClockTime: clock(),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := NewTotalGC(data.args.storage, data.args.options...)

			mock.AssertExpectationsForObjects(test, got.storage)
			assert.Equal(test, data.wantStorage, got.storage)

			// don't use the reflect.Value.Pointer() method for this check; see details:
			// * https://golang.org/pkg/reflect/#Value.Pointer
			// * https://stackoverflow.com/a/9644797
			require.NotNil(test, got.clock)
			assert.WithinDuration(test, data.wantClockTime, got.clock(), time.Hour)
		})
	}
}

func TestTotalGC_Clean(test *testing.T) {
	type fields struct {
		storage hashmap.Storage
		clock   cache.Clock
	}
	type args struct {
		ctx context.Context
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "without iterations",
			fields: fields{
				storage: func() hashmap.Storage {
					storage := new(MockStorage)
					storage.
						On("Iterate", mock.MatchedBy(func(handler hashmap.Handler) bool {
							return handler != nil
						})).
						Return(true).
						Once()

					return storage
				}(),
				clock: clock,
			},
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "with iterations",
			fields: fields{
				storage: func() hashmap.Storage {
					storage := new(MockStorage)
					storage.
						On("Iterate", mock.MatchedBy(func(handler hashmap.Handler) bool {
							return handler != nil
						})).
						Return(func(handler hashmap.Handler) bool {
							for i := 0; i < 15; i++ {
								var expirationTime time.Time
								if i < 3 {
									expirationTime = clock().Add(-time.Second)
								} else {
									expirationTime = clock().Add(time.Second)
								}

								if ok := handler(NewMockKeyWithID(23), cache.Value{
									Data:           "data",
									ExpirationTime: expirationTime,
								}); !ok {
									return false
								}
							}

							return true
						}).
						Once()
					storage.On("Delete", NewMockKeyWithID(23)).Times(3)

					return storage
				}(),
				clock: clock,
			},
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "with canceled iterations",
			fields: fields{
				storage: func() hashmap.Storage {
					storage := new(MockStorage)
					storage.
						On("Iterate", mock.MatchedBy(func(handler hashmap.Handler) bool {
							return handler != nil
						})).
						Return(func(handler hashmap.Handler) bool {
							for i := 0; i < 15; i++ {
								time.Sleep(timedTestDelay * 3 / 4)

								if ok := handler(NewMockKeyWithID(23), cache.Value{
									Data:           "data",
									ExpirationTime: clock().Add(-time.Second),
								}); !ok {
									return false
								}
							}

							return true
						}).
						Once()
					storage.On("Delete", NewMockKeyWithID(23)).Once()

					return storage
				}(),
				clock: clock,
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					go func() {
						time.Sleep(timedTestDelay)
						cancel()
					}()

					return ctx
				}(),
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gc := TotalGC{data.fields.storage, data.fields.clock}
			gc.Clean(data.args.ctx)

			mock.AssertExpectationsForObjects(test, data.fields.storage)
		})
	}
}

func TestTotalGC_handleIteration(test *testing.T) {
	type fields struct {
		storage hashmap.Storage
		clock   cache.Clock
	}
	type args struct {
		key   hashmap.Key
		value interface{}
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
		want   assert.BoolAssertionFunc
	}{
		{
			name: "with a not expired value",
			fields: fields{
				storage: new(MockStorage),
				clock:   clock,
			},
			args: args{
				key:   NewMockKeyWithID(23),
				value: cache.Value{Data: "data", ExpirationTime: clock().Add(time.Second)},
			},
			want: assert.True,
		},
		{
			name: "with an expired value",
			fields: fields{
				storage: func() hashmap.Storage {
					storage := new(MockStorage)
					storage.On("Delete", NewMockKeyWithID(23))

					return storage
				}(),
				clock: clock,
			},
			args: args{
				key:   NewMockKeyWithID(23),
				value: cache.Value{Data: "data", ExpirationTime: clock().Add(-time.Second)},
			},
			want: assert.True,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gc := TotalGC{data.fields.storage, data.fields.clock}
			got := gc.handleIteration(data.args.key, data.args.value)

			mock.AssertExpectationsForObjects(test, data.fields.storage, data.args.key)
			data.want(test, got)
		})
	}
}

func clock() time.Time {
	return time.Date(
		2006, time.January, 2, // year, month, day
		15, 4, 5, // hour, minute, second
		0,        // nanosecond
		time.UTC, // location
	)
}

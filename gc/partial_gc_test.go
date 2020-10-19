package gc

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cache "github.com/thewizardplusplus/go-cache"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

func TestNewPartialGC(test *testing.T) {
	storage := new(MockStorage)
	gc := NewPartialGC(storage, time.Now)

	mock.AssertExpectationsForObjects(test, storage)
	assert.Equal(test, storage, gc.storage)
	assert.Equal(test, getPointer(time.Now), getPointer(gc.clock))
}

// nolint: gocyclo
func TestPartialGC_Clean(test *testing.T) {
	type fields struct {
		storage Storage
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
				storage: func() Storage {
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
			name: "with a one try",
			fields: fields{
				storage: func() Storage {
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
			name: "with few tries",
			fields: fields{
				storage: func() Storage {
					var try int

					storage := new(MockStorage)
					storage.
						On("Iterate", mock.MatchedBy(func(handler hashmap.Handler) bool {
							return handler != nil
						})).
						Return(func(handler hashmap.Handler) bool {
							defer func() { try++ }()

							for i := 0; i < 15; i++ {
								var expirationTime time.Time
								if try == 0 && i < 5 {
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
						Twice()
					storage.On("Delete", NewMockKeyWithID(23)).Times(5)

					return storage
				}(),
				clock: clock,
			},
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "with canceled tries",
			fields: fields{
				storage: new(MockStorage),
				clock:   clock,
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()

					return ctx
				}(),
			},
		},
		{
			name: "with canceled iterations",
			fields: fields{
				storage: func() Storage {
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
			gc := PartialGC{data.fields.storage, data.fields.clock}
			gc.Clean(data.args.ctx)

			mock.AssertExpectationsForObjects(test, data.fields.storage)
		})
	}
}

func getPointer(value interface{}) uintptr {
	return reflect.ValueOf(value).Pointer()
}

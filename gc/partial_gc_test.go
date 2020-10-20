package gc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	cache "github.com/thewizardplusplus/go-cache"
	"github.com/thewizardplusplus/go-cache/models"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

func TestNewPartialGC(test *testing.T) {
	type args struct {
		storage hashmap.Storage
		options []PartialGCOption
	}

	for _, data := range []struct {
		name                  string
		args                  args
		wantStorage           hashmap.Storage
		wantClockTime         time.Time
		wantMaxIteratedCount  int
		wantMinExpiredPercent float64
	}{
		{
			name: "with default options",
			args: args{
				storage: new(MockStorage),
				options: nil,
			},
			wantStorage:           new(MockStorage),
			wantClockTime:         time.Now(),
			wantMaxIteratedCount:  defaultMaxIteratedCount,
			wantMinExpiredPercent: defaultMinExpiredPercent,
		},
		{
			name: "with the set clock",
			args: args{
				storage: new(MockStorage),
				options: []PartialGCOption{PartialGCWithClock(clock)},
			},
			wantStorage:           new(MockStorage),
			wantClockTime:         clock(),
			wantMaxIteratedCount:  defaultMaxIteratedCount,
			wantMinExpiredPercent: defaultMinExpiredPercent,
		},
		{
			name: "with the set maximal iterated count",
			args: args{
				storage: new(MockStorage),
				options: []PartialGCOption{PartialGCWithMaxIteratedCount(23)},
			},
			wantStorage:           new(MockStorage),
			wantClockTime:         time.Now(),
			wantMaxIteratedCount:  23,
			wantMinExpiredPercent: defaultMinExpiredPercent,
		},
		{
			name: "with the set minimal expired percent",
			args: args{
				storage: new(MockStorage),
				options: []PartialGCOption{PartialGCWithMinExpiredPercent(23)},
			},
			wantStorage:           new(MockStorage),
			wantClockTime:         time.Now(),
			wantMaxIteratedCount:  defaultMaxIteratedCount,
			wantMinExpiredPercent: 23,
		},
		{
			name: "with set options",
			args: args{
				storage: new(MockStorage),
				options: []PartialGCOption{
					PartialGCWithClock(clock),
					PartialGCWithMaxIteratedCount(23),
					PartialGCWithMinExpiredPercent(42),
				},
			},
			wantStorage:           new(MockStorage),
			wantClockTime:         clock(),
			wantMaxIteratedCount:  23,
			wantMinExpiredPercent: 42,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := NewPartialGC(data.args.storage, data.args.options...)

			mock.AssertExpectationsForObjects(test, got.storage)
			assert.Equal(test, data.wantStorage, got.storage)
			assert.Equal(test, data.wantMaxIteratedCount, got.maxIteratedCount)
			assert.Equal(test, data.wantMinExpiredPercent, got.minExpiredPercent)

			// don't use the reflect.Value.Pointer() method for this check; see details:
			// * https://golang.org/pkg/reflect/#Value.Pointer
			// * https://stackoverflow.com/a/9644797
			require.NotNil(test, got.clock)
			assert.WithinDuration(test, data.wantClockTime, got.clock(), time.Hour)
		})
	}
}

func TestPartialGC_Clean(test *testing.T) {
	type fields struct {
		storage           hashmap.Storage
		clock             models.Clock
		maxIteratedCount  int
		minExpiredPercent float64
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
				clock:             clock,
				maxIteratedCount:  20,
				minExpiredPercent: 0.25,
			},
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "with a one try",
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
				clock:             clock,
				maxIteratedCount:  20,
				minExpiredPercent: 0.25,
			},
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "with few tries",
			fields: fields{
				storage: func() hashmap.Storage {
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
				clock:             clock,
				maxIteratedCount:  20,
				minExpiredPercent: 0.25,
			},
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "with canceled tries",
			fields: fields{
				storage:           new(MockStorage),
				clock:             clock,
				maxIteratedCount:  20,
				minExpiredPercent: 0.25,
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
				clock:             clock,
				maxIteratedCount:  20,
				minExpiredPercent: 0.25,
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
			gc := PartialGC{
				storage:           data.fields.storage,
				clock:             data.fields.clock,
				maxIteratedCount:  data.fields.maxIteratedCount,
				minExpiredPercent: data.fields.minExpiredPercent,
			}
			gc.Clean(data.args.ctx)

			mock.AssertExpectationsForObjects(test, data.fields.storage)
		})
	}
}
